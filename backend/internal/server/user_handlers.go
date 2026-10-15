package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/z3vxo/TokenJack/internal/db"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed decoding json"})
		return
	}

	err := s.DB.CheckAuth(req.Usernmae, req.Password)
	switch {
	case err == nil:
		tokens, err := s.Auth.Generate(req.Usernmae)
		if err != nil {
			s.WriteJson(w, http.StatusInternalServerError, map[string]string{"Status": "error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
		return
	case errors.Is(err, db.ErrInvalidCredentials):
		s.WriteJson(w, http.StatusUnauthorized, map[string]string{"Status": "invalid creds"})

	default:
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"Status": "error"})
	}
}

func (s *Server) HandleMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := s.Auth.FromContext(r.Context())
	if !ok {
		s.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	s.WriteJson(w, http.StatusOK, map[string]string{"username": claims.Username})
}
