package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist
var dist embed.FS

var handler = newHandler()

func Handler() http.Handler { return handler }

func newHandler() http.Handler {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := adminPath(r.URL.Path)

		if !isFile(sub, p) {
			p = "/"
		}

		req := r.Clone(r.Context())
		req.URL.Path = p
		fileServer.ServeHTTP(w, req)
	})
}

func adminPath(urlPath string) string {
	p := strings.TrimPrefix(urlPath, "/admin")
	if p == "" {
		p = "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return path.Clean(p)
}

func isFile(fsys fs.FS, p string) bool {
	info, err := fs.Stat(fsys, strings.TrimPrefix(p, "/"))
	if err != nil {
		return false
	}
	return !info.IsDir()
}
