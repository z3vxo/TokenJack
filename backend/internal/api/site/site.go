package site

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/db"
	"github.com/z3vxo/TokenJack/internal/web"
)

type Handler struct {
	DB *db.DB
}

func New(d *db.DB) *Handler {
	return &Handler{DB: d}
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	web.Handler().ServeHTTP(w, r)
}

func (h *Handler) LandingPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "path")
	exits, err := h.DB.EnsureLandingIdExists(id)
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

func (h *Handler) LandingMisc(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "path")
	asset := chi.URLParam(r, "*")
	exits, err := h.DB.EnsureLandingIdExists(id)
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
}
