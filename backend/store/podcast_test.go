package store

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddPodcastToChannel(t *testing.T) {
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

			store, err := NewStore(testDir)
			assert.Nil(err)

			defer func() {
				store.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			cid, err := store.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), cid)

			if tt.wantErr {
				err := store.DropPodcasts()
				assert.Nil(err)
			}

			podcast, err := store.AddPodcastToChannel(cid, "podcast1.mp3", "podcast1", "10001")
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)
			assert.Equal(int64(1), podcast.ID)

			p := &Podcast{
				ID:    1,
				Title: "podcast1",
			}

			assert.Equal(p, podcast)
		})
	}
}

func TestUpdatingPodcast(t *testing.T) {
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

			store, err := NewStore(testDir)
			assert.Nil(err)

			defer func() {
				store.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			cid, err := store.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), cid)

			np, err := store.AddPodcastToChannel(cid, "podcast1.mp3", "podcast1", "10001")
			assert.Nil(err)
			assert.Equal(int64(1), np.ID)

			p, err := store.PodcastInfo(np.ID)
			assert.Nil(err)

			p.Title = "new title"
			p.Description = "short desc"
			p.Season = 1
			p.Episode = 3

			if tt.wantErr {
				err := store.DropPodcasts()
				assert.Nil(err)
			}

			err = store.UpdatePodcast(p)
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)

			pu, err := store.PodcastInfo(p.ID)
			assert.Nil(err)

			assert.Equal(int64(1), pu.ID)
			assert.Equal("podcast1.mp3", pu.Filename)
			assert.Equal(10001, pu.Length)
			assert.Equal("new title", pu.Title)
			assert.Equal("short desc", pu.Description)
			assert.Equal(1, pu.Season)
			assert.Equal(3, pu.Episode)
		})
	}
}

func TestListPodcasts(t *testing.T) {
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

			store, err := NewStore(testDir)
			assert.Nil(err)

			defer func() {
				store.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			cid, err := store.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), cid)

			for i := 1; i < 4; i++ {
				p, err := store.AddPodcastToChannel(cid, fmt.Sprintf("podcast%d.mp3", i), fmt.Sprintf("podcast%d", i), fmt.Sprintf("1000%d", i))
				assert.Nil(err)
				assert.Equal(int64(i), p.ID)
			}

			if tt.wantErr {
				err := store.DropPodcasts()
				assert.Nil(err)
			}

			ps, err := store.ListPodcastsFrom(cid)
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)
			assert.Equal(3, len(ps))
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

			store, err := NewStore(testDir)
			assert.Nil(err)

			defer func() {
				store.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			cid, err := store.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), cid)

			p, err := store.AddPodcastToChannel(cid, "podcast1.mp3", "podcast1", "10001")
			assert.Nil(err)
			assert.Equal(int64(1), p.ID)

			p2, err := store.AddPodcastToChannel(cid, "podcast2.mp3", "podcast2", "10002")
			assert.Nil(err)
			assert.Equal(int64(2), p2.ID)

			ps, err := store.ListPodcastsFrom(cid)
			assert.Nil(err)
			assert.Equal(2, len(ps))

			if tt.wantErr {
				err := store.DropPodcasts()
				assert.Nil(err)
			}

			err = store.DeletePodcast(p.ID)
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)

			ps, err = store.ListPodcastsFrom(cid)
			assert.Nil(err)
			assert.Equal(1, len(ps))
		})
	}
}
