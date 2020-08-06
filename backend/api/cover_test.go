package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/azzzak/fakecast/fs"
	"github.com/azzzak/fakecast/store"
	"github.com/stretchr/testify/assert"
)

func TestSetCoverURL(t *testing.T) {
	cfg := &Cfg{
		Host: "http://host.com",
	}

	type args struct {
		cfg *Cfg
		c   *store.Channel
	}
	tests := []struct {
		name string
		args args
		want *store.Channel
	}{
		{
			name: "empty",
			args: args{
				cfg: cfg,
				c:   &store.Channel{},
			},
			want: &store.Channel{},
		}, {
			name: "with cover",
			args: args{
				cfg: cfg,
				c: &store.Channel{
					Alias: "nope",
					Cover: "image.png",
				},
			},
			want: &store.Channel{
				Alias: "nope",
				Cover: fmt.Sprintf("%s/files/nope/%s/image.png", cfg.Host, fs.CoverDirName),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCoverURL(tt.args.cfg, tt.args.c)
			assert.Equal(t, tt.want, tt.args.c)
		})
	}
}

func TestUploadCover(t *testing.T) {
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

			path := filepath.Join("..", "_testdata", "tiny.jpg")
			file, err := os.Open(path)
			assert.Nil(err)

			fileContents, err := ioutil.ReadAll(file)
			assert.Nil(err)

			fi, err := file.Stat()
			assert.Nil(err)

			file.Close()

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", fi.Name())
			assert.Nil(err)

			part.Write(fileContents)

			err = writer.Close()
			assert.Nil(err)

			if tt.wantErr {
				err := s.DropChannels()
				assert.Nil(err)
			}

			r := httptest.NewRequest("POST", "/api/channel/1/cover/upload", body)
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

			path = filepath.Join(testDir, fs.PodcastsDirName, c.Alias, fs.CoverDirName, "tiny.jpg")
			file, err = os.Open(path)
			assert.Nil(err)

			fi, err = file.Stat()
			assert.Nil(err)

			file.Close()
		})
	}
}

func TestDeleteCover(t *testing.T) {
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

			coverName := "image.jpg"

			c.Alias = "1"
			c.Title = "channel 1"
			c.Cover = coverName

			err = s.UpdateChannel(c)
			assert.Nil(err)

			err = root.CreateDir(c.ID)
			assert.Nil(err)

			coverPath := filepath.Join(testDir, fs.PodcastsDirName, c.Alias, fs.CoverDirName, coverName)
			_, err = os.Create(coverPath)
			assert.Nil(err)

			if tt.wantErr {
				err := s.DropChannels()
				assert.Nil(err)
			}

			r := httptest.NewRequest("DELETE", fmt.Sprintf("/api/channel/1/cover/%s", coverName), nil)

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

			c, err = s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal("", c.Cover)

			_, err = os.Stat(coverPath)
			assert.NotNil(err)
		})
	}
}
