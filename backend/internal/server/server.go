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
	"github.com/z3vxo/TokenJack/internal/auth"
	"github.com/z3vxo/TokenJack/internal/codes"
	"github.com/z3vxo/TokenJack/internal/config"
	"github.com/z3vxo/TokenJack/internal/db"
)

type Server struct {
	DB    *db.DB
	Cfg   *config.Config
	Auth  *auth.Auth
	Codes *codes.Codes
}

func New() (*Server, error) {
	d, err := db.New()
	if err != nil {
		return nil, err
	}
	c, err := config.Get()
	if err != nil {
		return nil, err
	}

	codes := codes.New(d)
	return &Server{
		DB:    d,
		Cfg:   c,
		Auth:  auth.New([]byte(c.Auth.AccessSecret), []byte(c.Auth.RefreshSecret), c.Auth.AccessTTL, c.Auth.RefreshTTL),
		Codes: codes,
	}, nil
}

func (s *Server) SetupRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/admin", s.HandleAdminRoute)
	r.Get("/admin/*", s.HandleAdminRoute)

	r.Route("/rest", func(r chi.Router) {
		r.Post("/user/login", s.HandleLogin)
		r.Post("/code/generate", s.GenerateCode)

		r.Group(func(r chi.Router) {
			r.Use(s.Auth.AuthWare())
			r.Get("/user/me", s.HandleMe)

			r.Get("/campaigns/list", s.HandleListCampaigns)
			r.Post("/campaign/new", s.HandleNewCampaign)

		})
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
