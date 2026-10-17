package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/api"
	"github.com/z3vxo/TokenJack/internal/auth"
	"github.com/z3vxo/TokenJack/internal/db"
)

type Handler struct {
	DB   *db.DB
	Auth *auth.Auth
}

func New(d *db.DB, a *auth.Auth) *Handler {
	return &Handler{DB: d, Auth: a}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.Login)
	r.Group(func(r chi.Router) {
		r.Use(h.Auth.AuthWare())
		r.Get("/me", h.Me)
	})
	return r
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.Error(w, http.StatusInternalServerError, "failed decoding json")
		return
	}

	err := h.DB.CheckAuth(req.Username, req.Password)
	switch {
	case err == nil:
		tokens, err := h.Auth.Generate(req.Username)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, "error")
			return
		}
		api.WriteJSON(w, http.StatusOK, tokens)
	case errors.Is(err, db.ErrInvalidCredentials):
		api.Error(w, http.StatusUnauthorized, "invalid creds")
	default:
		api.Error(w, http.StatusInternalServerError, "error")
	}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.Auth.FromContext(r.Context())
	if !ok {
		api.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	api.WriteJSON(w, http.StatusOK, map[string]string{"username": claims.Username})
}
