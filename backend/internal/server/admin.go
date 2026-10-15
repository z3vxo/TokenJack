package server

import (
	"net/http"

	"github.com/z3vxo/TokenJack/internal/web"
)

func (s *Server) HandleAdminRoute(w http.ResponseWriter, r *http.Request) {
	web.Handler().ServeHTTP(w, r)
}
