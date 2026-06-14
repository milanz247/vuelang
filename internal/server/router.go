package server

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"

	"go-cloud-erp/routes"
)

// setupAPI mounts /api/v1 and hands off to routes/api.go.
func (s *Server) setupAPI() {
	api := s.router.Group("/api/v1")
	routes.Register(api, s.db)
}

// ── Production: serve embedded ui/dist ───────────────────────────────────────

func (s *Server) setupStaticProd(staticFS fs.FS) {
	publicFS, err := fs.Sub(staticFS, "ui/dist")
	if err != nil {
		panic("ui/dist not found in binary — run 'make build' first: " + err.Error())
	}

	s.router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" || path == "" {
			serveIndex(c, publicFS)
			return
		}
		f, err := publicFS.Open(path[1:])
		if err == nil {
			defer f.Close()
			if st, err := f.Stat(); err == nil && !st.IsDir() {
				c.FileFromFS(path[1:], http.FS(publicFS))
				return
			}
		}
		serveIndex(c, publicFS) // SPA fallback
	})
}

func serveIndex(c *gin.Context, publicFS fs.FS) {
	data, err := fs.ReadFile(publicFS, "index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "index.html missing")
		return
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

// ── Dev: reverse-proxy to Vite ────────────────────────────────────────────────

func (s *Server) setupStaticDev(viteAddr string) {
	target, err := url.Parse(viteAddr)
	if err != nil {
		panic("invalid vite URL: " + err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host   = target.Host
		req.Host        = target.Host
	}
	s.router.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
