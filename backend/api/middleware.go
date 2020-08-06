package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func (cfg *Cfg) channelID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := chi.URLParam(r, "channel")
		id, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			id = 0
		}
		ctx := context.WithValue(r.Context(), CID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *Cfg) podcastID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := chi.URLParam(r, "podcast")
		id, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			id = 0
		}
		ctx := context.WithValue(r.Context(), PID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func basicAuth(realm string, creds map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w, realm)
				return
			}

			credPass, credUserOk := creds[user]
			if !credUserOk || pass != credPass {
				basicAuthFailed(w, realm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}

func void(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
