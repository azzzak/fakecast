package store

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddChannel(t *testing.T) {
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

			if tt.wantErr {
				err := store.DropChannels()
				assert.Nil(err)
			}

			id, err := store.AddChannel()
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Equal(int64(1), id)
			assert.Nil(err)
		})
	}
}

func TestChannelInfo(t *testing.T) {
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

			id, err := store.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			if tt.wantErr {
				err := store.DropChannels()
				assert.Nil(err)
			}

			c, err := store.ChannelInfo(id)
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)
			assert.Equal(int64(1), c.ID)
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

			store, err := NewStore(testDir)
			assert.Nil(err)

			defer func() {
				store.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			id, err := store.AddChannel()
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c, err := store.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), id)

			c.Alias = "short1"
			c.Title = "title1"

			if tt.wantErr {
				err := store.DropChannels()
				assert.Nil(err)
			}

			err = store.UpdateChannel(c)
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)

			c, err = store.ChannelInfo(id)
			assert.Nil(err)
			assert.Equal(int64(1), c.ID)
			assert.Equal("short1", c.Alias)
			assert.Equal("title1", c.Title)
		})
	}
}

func TestListChannel(t *testing.T) {
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

			for i := 1; i < 4; i++ {
				id, err := store.AddChannel()
				assert.Nil(err)
				assert.Equal(int64(i), id)

				c, err := store.ChannelInfo(id)
				assert.Nil(err)
				assert.Equal(int64(i), c.ID)

				c.Alias = fmt.Sprintf("short%d", i)
				c.Title = fmt.Sprintf("title%d", i)

				err = store.UpdateChannel(c)
				assert.Nil(err)
			}

			if tt.wantErr {
				err := store.DropChannels()
				assert.Nil(err)
			}

			cs, err := store.ListChannels()
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)
			assert.Equal(3, len(cs))
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

			store, err := NewStore(testDir)
			assert.Nil(err)

			defer func() {
				store.Close()
				err = os.RemoveAll(testDir)
				assert.Nil(err)
			}()

			for i := 1; i < 4; i++ {
				id, err := store.AddChannel()
				assert.Nil(err)
				assert.Equal(int64(i), id)

				c, err := store.ChannelInfo(id)
				assert.Nil(err)
				assert.Equal(int64(i), c.ID)

				c.Alias = fmt.Sprintf("short%d", i)
				c.Title = fmt.Sprintf("title%d", i)

				err = store.UpdateChannel(c)
				assert.Nil(err)
			}

			cs, err := store.ListChannels()
			if err != nil {
				t.Fatal("error while listing channels")
			}

			assert.Equal(3, len(cs))

			if tt.wantErr {
				err := store.DropChannels()
				assert.Nil(err)
			}

			err = store.DeleteChannel(2)
			if tt.wantErr {
				assert.NotNil(err)
				return
			}

			assert.Nil(err)

			cs, err = store.ListChannels()
			assert.Nil(err)

			assert.Equal(2, len(cs))
		})
	}
}
