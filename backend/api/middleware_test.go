package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestChannelID(t *testing.T) {
	assert := assert.New(t)
	cfg := &Cfg{}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := r.Context().Value(CID)
		assert.NotNil(cid)
		cidInt, ok := cid.(int64)
		assert.Equal(true, ok)
		assert.Equal(int64(123), cidInt)
	})

	r := httptest.NewRequest("GET", "/", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("channel", "123")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	handlerToTest := cfg.channelID(nextHandler)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), r)
}

func TestPodcastID(t *testing.T) {
	assert := assert.New(t)
	cfg := &Cfg{}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pid := r.Context().Value(PID)
		assert.NotNil(pid)
		pidInt, ok := pid.(int64)
		assert.Equal(true, ok)
		assert.Equal(int64(321), pidInt)
	})

	r := httptest.NewRequest("GET", "/", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("podcast", "321")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	handlerToTest := cfg.podcastID(nextHandler)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), r)
}
