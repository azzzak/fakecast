package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
	"github.com/go-chi/chi"
)

type ctxKey int

const (
	// CID id of a channel
	CID ctxKey = iota + 1
	// PID id of a podcast
	PID
)

// Cfg configuration
type Cfg struct {
	Store      store.Store
	FS         fs.Dir
	Host       string
	Credential string
}

type hndlr func(http.ResponseWriter, *http.Request) error

func (fn hndlr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)

		switch err.(type) {
		case *store.Error:
			fmt.Fprintf(os.Stderr, "Database error: %v\n", err)
		case *fs.Error:
			fmt.Fprintf(os.Stderr, "Filesystem error: %v\n", err)
		default:
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

func credential(cred string) map[string]string {
	m := map[string]string{}
	c := strings.SplitN(cred, ":", 2)

	switch len(c) {
	case 2:
		m[c[0]] = c[1]
	case 1:
		m["fakecast"] = c[0]
	}

	return m
}

// InitHandlers for API
func InitHandlers(cfg *Cfg) *chi.Mux {
	var baseURL string

	base, err := url.Parse(cfg.Host)
	if err != nil {
		fmt.Println("HOST env variable is not valid URL")
		os.Exit(1)
	}

	if base.Path != "" {
		baseURL = base.Path
	}

	auth := void
	if cfg.Credential != "" {
		auth = basicAuth("auth", credential(cfg.Credential))
	}

	r := chi.NewRouter()

	r.Use(corsMiddleware().Handler)

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Cant't get working directory")
		os.Exit(1)
	}

	guiDir := protect{
		http.Dir(filepath.Join(workDir, fs.FrontDirName)),
	}

	podcastsDir := protect{
		fs: http.Dir(cfg.FS.Root),
	}

	fileServer(r, "/"+strings.TrimPrefix(baseURL, "/"), guiDir)

	fileServer(r, baseURL+"/files", podcastsDir)

	r.Get(baseURL+"/feed/{channel}", hndlr(cfg.genFeed).ServeHTTP)

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		robots := `User-agent: *
Disallow: /`
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(robots))
	})

	r.With(auth).Route(baseURL+"/api", func(r chi.Router) {
		r.Get("/list", hndlr(cfg.list).ServeHTTP)

		r.Route("/channel", func(r chi.Router) {
			r.Post("/", hndlr(cfg.createChannel).ServeHTTP)

			r.With(cfg.channelID).Route("/{channel}", func(r chi.Router) {
				r.Get("/", hndlr(cfg.overview).ServeHTTP)
				r.Put("/", hndlr(cfg.updateChannel).ServeHTTP)
				r.Delete("/", hndlr(cfg.deleteChannel).ServeHTTP)
				r.Post("/upload", hndlr(cfg.uploadPodcast).ServeHTTP)

				r.Post("/cover/upload", hndlr(cfg.uploadCover).ServeHTTP)
				r.Delete("/cover/{cover}", hndlr(cfg.deleteCover).ServeHTTP)

				r.Route("/podcast", func(r chi.Router) {
					r.With(cfg.podcastID).Route("/{podcast}", func(r chi.Router) {
						r.Get("/", hndlr(cfg.podcastInfo).ServeHTTP)
						r.Put("/", hndlr(cfg.updatePodcast).ServeHTTP)
						r.Delete("/", hndlr(cfg.deletePodcast).ServeHTTP)
					})
				})
			})
		})
	})

	return r
}
