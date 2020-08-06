package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestCheckChannels(t *testing.T) {
	assert := assert.New(t)
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

	for i := 1; i < 4; i++ {
		id, err := s.AddChannel()
		assert.Nil(err)
		assert.Equal(int64(i), id)

		c, err := s.ChannelInfo(id)
		assert.Nil(err)
		assert.Equal(int64(i), c.ID)

		c.Alias = fmt.Sprintf("%d", i)
		c.Title = fmt.Sprintf("New channel %d", i)

		err = s.UpdateChannel(c)
		assert.Nil(err)

		err = root.CreateDir(c.ID)
		assert.Nil(err)
	}

	cs, err := s.ListChannels()
	assert.Nil(err)

	type args struct {
		cfg *Cfg
		cs  []store.Channel
	}
	tests := []struct {
		name   string
		args   args
		delete []string
		want   []store.Channel
	}{
		{
			name: "full",
			args: args{
				cfg: cfg,
				cs:  cs,
			},
			delete: []string{},
			want:   cs,
		}, {
			name: "remove one",
			args: args{
				cfg: cfg,
				cs:  cs,
			},
			delete: []string{"2"},
			want: []store.Channel{
				{ID: 1, Alias: "1", Title: "New channel 1"},
				{ID: 3, Alias: "3", Title: "New channel 3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, del := range tt.delete {
				path := filepath.Join(cfg.FS.Root, del)
				err := os.RemoveAll(path)
				assert.Nil(err)
			}

			got := checkChannels(tt.args.cfg, tt.args.cs)
			assert.Equal(tt.want, got)
		})
	}
}
func TestCheckPodcasts(t *testing.T) {
	assert := assert.New(t)
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

	ps := []store.Podcast{
		{ID: 3, Filename: "podcast3.mp3", Title: "podcast3"},
		{ID: 2, Filename: "podcast2.mp3", Title: "podcast2"},
		{ID: 1, Filename: "podcast1.mp3", Title: "podcast1"},
	}

	type args struct {
		cfg   *Cfg
		alias string
		ps    []store.Podcast
	}
	tests := []struct {
		name   string
		args   args
		delete []string
		want   []store.Podcast
	}{
		{
			name: "full",
			args: args{
				cfg:   cfg,
				alias: "1",
				ps:    append([]store.Podcast(nil), ps...),
			},
			delete: []string{},
			want:   append([]store.Podcast(nil), ps...),
		}, {
			name: "remove one",
			args: args{
				cfg:   cfg,
				alias: "1",
				ps:    append([]store.Podcast(nil), ps...),
			},
			delete: []string{"podcast2.mp3"},
			want: []store.Podcast{
				{ID: 3, Filename: "podcast3.mp3", Title: "podcast3"},
				{ID: 1, Filename: "podcast1.mp3", Title: "podcast1"},
			},
		}, {
			name: "remove two",
			args: args{
				cfg:   cfg,
				alias: "1",
				ps:    append([]store.Podcast(nil), ps...),
			},
			delete: []string{"podcast1.mp3", "podcast2.mp3"},
			want: []store.Podcast{
				{ID: 3, Filename: "podcast3.mp3", Title: "podcast3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := s.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)
			defer s.DeleteChannel(id)

			c, err := s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)

			c.Alias = "1"
			c.Title = "New channel 1"

			err = s.UpdateChannel(c)
			assert.Nil(err)

			err = root.CreateDir(c.ID)
			assert.Nil(err)
			defer os.RemoveAll(filepath.Join(root.Root, c.Alias))

			for i := 1; i < 4; i++ {
				_, err := s.AddPodcastToChannel(c.ID, fmt.Sprintf("podcast%d.mp3", i), fmt.Sprintf("podcast%d", i), fmt.Sprintf("1000%d", i))
				assert.Nil(err)

				_, err = os.Create(filepath.Join(testDir, fs.PodcastsDirName, c.Alias, fmt.Sprintf("podcast%d.mp3", i)))
				assert.Nil(err)
			}

			for _, del := range tt.delete {
				path := filepath.Join(cfg.FS.Root, tt.args.alias, del)
				err := os.Remove(path)
				assert.Nil(err)
			}

			got := checkPodcasts(tt.args.cfg, tt.args.alias, tt.args.ps)
			assert.Equal(tt.want, got)
		})
	}
}

func TestCreateChannel(t *testing.T) {
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
			cfg := &Cfg{
				Store: s,
				FS:    fs.NewRoot(testDir),
			}

			defer func() {
				s.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			if tt.wantErr {
				err := s.DropChannels()
				assert.Nil(err)
			}

			for i := 1; i < 4; i++ {
				r := httptest.NewRequest("POST", "/api/channel", nil)

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

				var c store.Channel

				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&c)
				assert.Nil(err)

				tc := store.Channel{
					ID:    int64(i),
					Alias: strconv.Itoa(i),
					Title: fmt.Sprintf("New channel %d", i),
				}

				assert.Equal(tc, c)
			}
		})
	}
}

func TestList(t *testing.T) {
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

			fs := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    fs,
			}

			defer func() {
				s.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			for i := 1; i < 4; i++ {
				id, err := s.AddChannel()
				assert.Nil(err)
				assert.Equal(int64(i), id)

				c, err := s.ChannelInfo(id)
				assert.Nil(err)
				assert.Equal(int64(i), c.ID)

				c.Alias = fmt.Sprintf("%d", i)
				c.Title = fmt.Sprintf("New channel %d", i)

				err = s.UpdateChannel(c)
				assert.Nil(err)

				err = fs.CreateDir(c.ID)
				assert.Nil(err)
			}

			if tt.wantErr {
				err := s.DropChannels()
				assert.Nil(err)
			}

			r := httptest.NewRequest("GET", "/api/list", nil)

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

			var cs []store.Channel

			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&cs)
			assert.Nil(err)

			ts := []store.Channel{
				{ID: 1, Alias: "1", Title: "New channel 1"},
				{ID: 2, Alias: "2", Title: "New channel 2"},
				{ID: 3, Alias: "3", Title: "New channel 3"},
			}

			assert.Equal(ts, cs)
		})
	}
}
func TestOverview(t *testing.T) {
	assert := assert.New(t)
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
	c.Title = "channel 1"

	err = s.UpdateChannel(c)
	assert.Nil(err)

	err = root.CreateDir(c.ID)
	assert.Nil(err)

	tests := []struct {
		name        string
		n           int
		want        overview
		errChannels bool
		errPodcasts bool
	}{
		{
			name: "w/o podcasts",
			n:    0,
			want: overview{
				Channel: &store.Channel{ID: 1, Alias: c.Alias, Title: c.Title},
			},
			errChannels: false,
			errPodcasts: false,
		}, {
			name: "with podcasts",
			n:    3,
			want: overview{
				Channel: &store.Channel{ID: 1, Alias: c.Alias, Title: c.Title},
				Podcasts: []store.Podcast{
					{
						ID:       3,
						Filename: "podcast3.mp3",
						Title:    "podcast3",
					}, {
						ID:       2,
						Filename: "podcast2.mp3",
						Title:    "podcast2",
					}, {
						ID:       1,
						Filename: "podcast1.mp3",
						Title:    "podcast1",
					},
				},
			},
			errChannels: false,
			errPodcasts: false,
		}, {
			name:        "error channels",
			n:           0,
			want:        overview{},
			errChannels: true,
			errPodcasts: false,
		}, {
			name:        "error podcasts",
			n:           0,
			want:        overview{},
			errChannels: false,
			errPodcasts: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for i := 1; i <= tt.n; i++ {
				_, err := s.AddPodcastToChannel(c.ID, fmt.Sprintf("podcast%d.mp3", i), fmt.Sprintf("podcast%d", i), fmt.Sprintf("1000%d", i))
				assert.Nil(err)

				_, err = os.Create(filepath.Join(testDir, fs.PodcastsDirName, c.Alias, fmt.Sprintf("podcast%d.mp3", i)))
				assert.Nil(err)
			}

			if tt.errChannels {
				err := s.DropChannels()
				assert.Nil(err)
			}

			if tt.errPodcasts {
				err := s.DropPodcasts()
				assert.Nil(err)
			}

			r := httptest.NewRequest("GET", "/api/channel/1", nil)

			w := httptest.NewRecorder()
			handler := http.Handler(InitHandlers(cfg))
			handler.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.errChannels || tt.errPodcasts {
				assert.Equal(http.StatusInternalServerError, resp.StatusCode)
				return
			}

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				t.Fatalf("Got status code: %d\n", resp.StatusCode)
			}

			var o overview

			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&o)
			assert.Nil(err)

			assert.Equal(tt.want, o)
		})
	}
}

func TestUpdateChannel(t *testing.T) {
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

			fs := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    fs,
			}

			id, err := s.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c, err := s.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)

			c.Alias = "update alias"
			c.Title = "update channel"
			c.Description = "new desc"
			c.Cover = "cover.png"
			c.Author = "user"

			err = fs.CreateDir(c.ID)
			assert.Nil(err)

			u := updateChannel{
				Channel:  c,
				OldAlias: "1",
			}

			jsonStr, err := json.Marshal(u)
			assert.Nil(err)

			if tt.wantErr {
				err := s.DropChannels()
				assert.Nil(err)
			}

			r := httptest.NewRequest("PUT", "/api/channel/1", bytes.NewBuffer(jsonStr))

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

			var ru updateResponse

			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&ru)
			assert.Nil(err)

			want := &store.Channel{
				ID:          int64(1),
				Alias:       "update alias",
				Title:       "update channel",
				Description: "new desc",
				Cover:       "cover.png",
				Author:      "user",
			}

			info, err := s.ChannelInfo(1)
			assert.Nil(err)
			assert.Equal(want, info)
		})
	}
}

func TestDeleteChannel(t *testing.T) {
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

			fs := fs.NewRoot(testDir)

			cfg := &Cfg{
				Store: s,
				FS:    fs,
			}

			for i := 1; i < 4; i++ {
				id, err := s.AddChannel()
				assert.Nil(err)
				assert.Equal(int64(i), id)

				c, err := s.ChannelInfo(id)
				assert.Nil(err)
				assert.Equal(int64(i), c.ID)

				c.Alias = fmt.Sprintf("%d", i)
				c.Title = fmt.Sprintf("New channel %d", i)

				err = s.UpdateChannel(c)
				assert.Nil(err)

				err = fs.CreateDir(c.ID)
				assert.Nil(err)
			}

			if tt.wantErr {
				err := s.DropChannels()
				assert.Nil(err)
			}

			r := httptest.NewRequest("DELETE", "/api/channel/2", nil)

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

			cs, err := s.ListChannels()
			assert.Nil(err)

			want := []store.Channel{
				{ID: 1, Alias: "1", Title: "New channel 1"},
				{ID: 3, Alias: "3", Title: "New channel 3"},
			}

			assert.Equal(want, cs)
		})
	}
}
