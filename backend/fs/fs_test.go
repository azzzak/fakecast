package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRoot(t *testing.T) {
	tests := []struct {
		name string
		root string
		want Dir
	}{
		{
			name: "default",
			root: "/home/fakecast",
			want: Dir{
				Root: fmt.Sprintf("/home/fakecast/%s", PodcastsDirName),
			},
		}, {
			name: "empty",
			root: "",
			want: Dir{
				Root: fmt.Sprintf("%s", PodcastsDirName),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRoot(tt.root)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateDir(t *testing.T) {
	assert := assert.New(t)
	testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())
	defer func() {
		err := os.RemoveAll(testDir)
		assert.Nil(err)
	}()

	d := &Dir{
		Root: filepath.Join(testDir, PodcastsDirName),
	}

	channel := int64(1)
	channelStr := strconv.FormatInt(channel, 10)

	err := d.CreateDir(channel)
	assert.Nil(err)

	path := filepath.Join(d.Root, channelStr, CoverDirName)
	_, err = os.Stat(path)
	assert.Nil(err)

	err = d.CreateDir(channel)
	assert.Nil(err)
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())
	defer func() {
		err := os.RemoveAll(testDir)
		assert.Nil(err)
	}()

	d := &Dir{
		Root: filepath.Join(testDir, PodcastsDirName),
	}

	channel := int64(1)
	channelStr := strconv.FormatInt(channel, 10)
	err := d.CreateDir(channel)
	assert.Nil(err)

	path := filepath.Join(d.Root, channelStr, CoverDirName)
	_, err = os.Stat(path)
	assert.Nil(err)

	coverName := "image.jpg"
	coverPath := filepath.Join(d.Root, channelStr, CoverDirName, coverName)
	_, err = os.Create(coverPath)
	assert.Nil(err)

	podcastName := "podcast.mp3"
	podcastPath := filepath.Join(d.Root, channelStr, podcastName)
	_, err = os.Create(podcastPath)
	assert.Nil(err)

	err = d.RemoveCover(channelStr, coverName)
	assert.Nil(err)
	_, err = os.Stat(coverPath)
	assert.NotNil(err)
	err = d.RemoveCover(channelStr, coverName)
	assert.NotNil(err)

	err = d.RemovePodcast(channelStr, podcastName)
	assert.Nil(err)
	_, err = os.Stat(podcastPath)
	assert.NotNil(err)
	err = d.RemovePodcast(channelStr, podcastName)
	assert.NotNil(err)

	err = d.RemoveDir(channelStr)
	assert.Nil(err)
	_, err = os.Stat(path)
	assert.NotNil(err)
}

func TestSave(t *testing.T) {
	assert := assert.New(t)
	testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())
	defer func() {
		err := os.RemoveAll(testDir)
		assert.Nil(err)
	}()

	d := &Dir{
		Root: filepath.Join(testDir, PodcastsDirName),
	}

	channel := int64(1)
	channelStr := strconv.FormatInt(channel, 10)
	err := d.CreateDir(channel)
	assert.Nil(err)

	filename := "data.txt"
	data := "123"

	p, err := d.SavePodcastToDir(channelStr, filename)
	assert.Nil(err)
	_, err = p.WriteString(data)
	assert.Nil(err)
	err = p.Close()
	assert.Nil(err)

	path := filepath.Join(d.Root, channelStr, filename)
	got, err := ioutil.ReadFile(path)
	assert.Nil(err)
	assert.Equal(data, string(got))

	err = os.Remove(path)
	assert.Nil(err)

	c, err := d.SaveCover(channelStr, filename)
	assert.Nil(err)
	_, err = c.WriteString(data)
	assert.Nil(err)
	err = c.Close()
	assert.Nil(err)

	path = filepath.Join(d.Root, channelStr, CoverDirName, filename)
	got, err = ioutil.ReadFile(path)
	assert.Nil(err)
	assert.Equal(data, string(got))
}

func TestRenameDir(t *testing.T) {
	assert := assert.New(t)
	testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())
	defer func() {
		err := os.RemoveAll(testDir)
		assert.Nil(err)
	}()

	d := &Dir{
		Root: filepath.Join(testDir, PodcastsDirName),
	}

	c := int64(1)
	old := strconv.FormatInt(c, 10)
	new := "2"

	err := d.CreateDir(c)
	assert.Nil(err)

	path := filepath.Join(d.Root, old)
	_, err = os.Stat(path)
	assert.Nil(err)

	err = d.RenameDir(old, new)
	assert.Nil(err)

	path = filepath.Join(d.Root, old)
	_, err = os.Stat(path)
	assert.NotNil(err)

	path = filepath.Join(d.Root, new)
	_, err = os.Stat(path)
	assert.Nil(err)

	err = d.RenameDir(new, new)
	assert.NotNil(err)
}

func TestIsExist(t *testing.T) {
	assert := assert.New(t)
	testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())
	defer func() {
		err := os.RemoveAll(testDir)
		assert.Nil(err)
	}()

	d := &Dir{
		Root: filepath.Join(testDir, PodcastsDirName),
	}

	channel := int64(1)
	channelStr := strconv.FormatInt(channel, 10)
	err := d.CreateDir(channel)
	assert.Nil(err)

	path := filepath.Join(d.Root, channelStr)
	_, err = os.Stat(path)
	assert.Nil(err)

	podcastName := "podcast.mp3"
	podcastPath := filepath.Join(d.Root, channelStr, podcastName)
	_, err = os.Create(podcastPath)
	assert.Nil(err)

	pb := d.IsPodcastExist(channelStr, podcastName)
	assert.Equal(true, pb)
	err = os.Remove(podcastPath)
	assert.Nil(err)
	pb = d.IsPodcastExist(channelStr, podcastName)
	assert.Equal(false, pb)

	cb := d.IsDirExist(channelStr)
	assert.Equal(true, cb)
	err = os.RemoveAll(path)
	assert.Nil(err)
	cb = d.IsDirExist(channelStr)
	assert.Equal(false, cb)
}

func TestNameAndExtFrom(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantName string
		wantExt  string
	}{
		{
			name:     "simple",
			filename: "podcast.mp3",
			wantName: "podcast",
			wantExt:  "mp3",
		}, {
			name:     "simple mp3",
			filename: "podcast-ep2.mp3",
			wantName: "podcast-ep2",
			wantExt:  "mp3",
		}, {
			name:     "simple m4a",
			filename: "podcast-ep5.m4a",
			wantName: "podcast-ep5",
			wantExt:  "m4a",
		}, {
			name:     "two extensions",
			filename: "podcast-s2.1.mp3",
			wantName: "podcast-s2.1",
			wantExt:  "mp3",
		}, {
			name:     "without extension",
			filename: "podcast-ep4",
			wantName: "podcast-ep4",
			wantExt:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, ext := NameAndExtFrom(tt.filename)
			assert.Equal(t, tt.wantName, name)
			assert.Equal(t, tt.wantExt, ext)
		})
	}
}
