package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

func (s *Server) HandleLandingPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "path")
	exits, err := s.DB.EnsureLandingIdExists(id)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(exits)

	if exits {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	path := filepath.Join(home, ".local", "share", "token-jack", "user-files", id, "index.html")
	fmt.Println(path)
	http.ServeFile(w, r, path)
}

func (s *Server) HandleLandingMisc(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "path")
	asset := chi.URLParam(r, "*")
	exits, err := s.DB.EnsureLandingIdExists(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if exits {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	baseDir := filepath.Join(home, ".local", "share", "token-jack", "user-files", id)
	filePath := filepath.Join(baseDir, filepath.Clean(asset))

	rel, err := filepath.Rel(baseDir, filePath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	fmt.Println(filePath)

	http.ServeFile(w, r, filePath)
	return

}
