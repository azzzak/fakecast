package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
	"github.com/stretchr/testify/assert"
)

func TestUploadPodcast(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		}, {
			name:    "error",
			wantErr: true,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())

			err := os.MkdirAll(testDir, os.ModePerm)
			assert.Nil(err)

			s, err := store.NewStore(testDir)
			assert.Nil(err)
			defer func() {
				s.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			root := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    root,
			}

			id, err := s.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c, err := s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)

			c.Alias = "1"

			err = s.UpdateChannel(c)
			assert.Nil(err)

			err = root.CreateDir(c.ID)
			assert.Nil(err)

			path := filepath.Join("..", "_testdata", "tiny.mp3")
			file, err := os.Open(path)
			assert.Nil(err)

			fileContents, err := ioutil.ReadAll(file)
			assert.Nil(err)
			length := 72

			fi, err := file.Stat()
			assert.Nil(err)

			file.Close()

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", fi.Name())
			assert.Nil(err)

			part.Write(fileContents)

			writer.WriteField("length", strconv.Itoa(length))

			err = writer.Close()
			assert.Nil(err)

			if tt.wantErr {
				err := s.DropPodcasts()
				assert.Nil(err)
			}

			r := httptest.NewRequest("POST", "/api/channel/1/upload", body)
			r.Header.Add("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			handler := http.Handler(InitHandlers(cfg))
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.wantErr {
				assert.Equal(http.StatusInternalServerError, resp.StatusCode)
				return
			}

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				t.Fatalf("Got status code: %d\n", resp.StatusCode)
			}

			path = filepath.Join(testDir, fs.PodcastsDirName, c.Alias, "tiny.mp3")
			file, err = os.Open(path)
			assert.Nil(err)

			fi, err = file.Stat()
			assert.Nil(err)

			file.Close()

			assert.Equal(int64(length), fi.Size())
		})
	}
}

func TestPodcastInfo(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		}, {
			name:    "error",
			wantErr: true,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())

			err := os.MkdirAll(testDir, os.ModePerm)
			assert.Nil(err)

			s, err := store.NewStore(testDir)
			assert.Nil(err)
			defer func() {
				s.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			root := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    root,
			}

			id, err := s.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c, err := s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)

			c.Alias = "1"

			err = s.UpdateChannel(c)
			assert.Nil(err)

			err = root.CreateDir(c.ID)
			assert.Nil(err)

			_, err = s.AddPodcastToChannel(c.ID, "podcast.mp3", "podcast", "10001")
			assert.Nil(err)

			_, err = os.Create(filepath.Join(testDir, fs.PodcastsDirName, "1", "podcast.mp3"))
			assert.Nil(err)

			if tt.wantErr {
				err := s.DropPodcasts()
				assert.Nil(err)
			}

			r := httptest.NewRequest("GET", "/api/channel/1/podcast/1", nil)

			w := httptest.NewRecorder()
			handler := http.Handler(InitHandlers(cfg))
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.wantErr {
				assert.Equal(http.StatusInternalServerError, resp.StatusCode)
				return
			}

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				t.Fatalf("Got status code: %d\n", resp.StatusCode)
			}

			var p store.Podcast

			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&p)
			assert.Nil(err)

			want := store.Podcast{
				ID:       1,
				Filename: "podcast.mp3",
				Title:    "podcast",
				Length:   10001,
			}

			assert.Equal(want, p)
		})
	}
}

func TestUpdatePodcast(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		}, {
			name:    "error",
			wantErr: true,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())

			err := os.MkdirAll(testDir, os.ModePerm)
			assert.Nil(err)

			s, err := store.NewStore(testDir)
			assert.Nil(err)
			defer func() {
				s.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			root := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    root,
			}

			id, err := s.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c, err := s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)

			c.Alias = "1"

			err = s.UpdateChannel(c)
			assert.Nil(err)

			err = root.CreateDir(c.ID)
			assert.Nil(err)

			_, err = s.AddPodcastToChannel(c.ID, "podcast.mp3", "p", "10001")
			assert.Nil(err)

			p, err := s.PodcastInfo(1)
			assert.Nil(err)

			p.Title = "podcast"
			p.Description = "desc"
			p.Duration = 101
			p.Artwork = "omahe.png"
			p.Explicit = 1
			p.Season = 2
			p.Episode = 3

			jsonStr, err := json.Marshal(p)
			assert.Nil(err)

			if tt.wantErr {
				err := s.DropPodcasts()
				assert.Nil(err)
			}

			r := httptest.NewRequest("PUT", "/api/channel/1/podcast/1", bytes.NewBuffer(jsonStr))

			w := httptest.NewRecorder()
			handler := http.Handler(InitHandlers(cfg))
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.wantErr {
				assert.Equal(http.StatusInternalServerError, resp.StatusCode)
				return
			}

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				t.Fatalf("Got status code: %d\n", resp.StatusCode)
			}

			want := &store.Podcast{
				ID:          int64(1),
				Filename:    "podcast.mp3",
				Length:      10001,
				Title:       "podcast",
				Description: "desc",
				Duration:    101,
				Artwork:     "omahe.png",
				Explicit:    1,
				Season:      2,
				Episode:     3,
			}

			info, err := s.PodcastInfo(1)
			assert.Nil(err)

			info.GUID = ""
			info.PubDate = ""

			assert.Equal(want, info)
		})
	}
}

func TestDeletePodcast(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		}, {
			name:    "error",
			wantErr: true,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())

			err := os.MkdirAll(testDir, os.ModePerm)
			assert.Nil(err)

			s, err := store.NewStore(testDir)
			assert.Nil(err)
			defer func() {
				s.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			root := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    root,
			}

			id, err := s.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c, err := s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)

			c.Alias = "1"

			err = s.UpdateChannel(c)
			assert.Nil(err)

			err = root.CreateDir(c.ID)
			assert.Nil(err)

			for i := 1; i < 4; i++ {
				_, err := s.AddPodcastToChannel(c.ID, fmt.Sprintf("podcast%d.mp3", i), fmt.Sprintf("podcast%d", i), fmt.Sprintf("1000%d", i))
				assert.Nil(err)

				_, err = os.Create(filepath.Join(testDir, fs.PodcastsDirName, "1", fmt.Sprintf("podcast%d.mp3", i)))
				assert.Nil(err)
			}

			if tt.wantErr {
				err := s.DropPodcasts()
				assert.Nil(err)
			}

			r := httptest.NewRequest("DELETE", "/api/channel/1/podcast/2", nil)

			w := httptest.NewRecorder()
			handler := http.Handler(InitHandlers(cfg))
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.wantErr {
				assert.Equal(http.StatusInternalServerError, resp.StatusCode)
				return
			}

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				t.Fatalf("Got status code: %d\n", resp.StatusCode)
			}

			ps, err := s.ListFullPodcastsFrom(c.ID)
			assert.Nil(err)

			want := []store.Podcast{
				{ID: 3, Filename: "podcast3.mp3", Title: "podcast3", Length: 10003},
				{ID: 1, Filename: "podcast1.mp3", Title: "podcast1", Length: 10001},
			}

			assert.Equal(want, ps)
		})
	}
}
