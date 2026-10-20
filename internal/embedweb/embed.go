package embedweb

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//go:embed dist
var files embed.FS

// Mode controls how the frontend is served.
type Mode string

const (
	ModeEmbedded Mode = "embedded" // serve from embedded dist/ (default)
	ModeExternal Mode = "external" // proxy to an external dev server
)

var (
	currentMode  Mode = ModeEmbedded
	externalProxy *httputil.ReverseProxy
)

// SetMode configures the frontend serving mode.
// For external mode, proxyURL is the external frontend server (e.g. "http://localhost:3001").
func SetMode(mode Mode, proxyURL string) {
	currentMode = mode
	if mode == ModeExternal && proxyURL != "" {
		if u, err := url.Parse(proxyURL); err == nil {
			externalProxy = httputil.NewSingleHostReverseProxy(u)
		}
	}
}

// Handler returns an http.Handler that serves the frontend.
// In embedded mode it serves static files + SPA fallback from the embedded dist/.
// In external mode it proxies all requests to the external dev server.
func Handler() http.Handler {
	if currentMode == ModeExternal && externalProxy != nil {
		return externalProxy
	}
	return embeddedHandler()
}

func embeddedHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.FileServer(AssetsFS()))
	indexHTML, err := ReadFile("index.html")
	if err != nil {
		log.Fatalf("embedweb: missing dist/index.html – run scripts/copy-web-dist.sh before building")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	return mux
}

// Dist returns the frontend dist filesystem (rooted at dist/).
func Dist() fs.FS {
	dist, err := fs.Sub(files, "dist")
	if err != nil {
		log.Fatalf("embedweb: %v", err)
	}
	return dist
}

// FS returns the frontend dist as http.FileSystem (rooted at dist/).
func FS() http.FileSystem {
	return http.FS(Dist())
}

// AssetsFS returns a filesystem rooted at dist/assets for serving bundled
// JS/CSS/image assets. Use with: r.StaticFS("/assets", embedweb.AssetsFS())
func AssetsFS() http.FileSystem {
	assets, err := fs.Sub(files, "dist/assets")
	if err != nil {
		log.Fatalf("embedweb assets: %v", err)
	}
	return http.FS(assets)
}

// ReadFile is a convenience wrapper around fs.ReadFile on the dist filesystem.
func ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(Dist(), name)
}