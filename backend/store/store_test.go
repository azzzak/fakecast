package store

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	assert := assert.New(t)
	testDir := fmt.Sprintf("test_dir_%x", time.Now().Unix())

	err := os.MkdirAll(testDir, os.ModePerm)
	assert.Nil(err)

	s, err := NewStore(testDir)
	assert.Nil(err)

	defer func() {
		s.Close()
		err = os.RemoveAll(testDir)
		assert.Nil(err)
	}()

	row := s.db.QueryRow("SELECT count(*) FROM channels")

	var n int
	err = row.Scan(&n)
	assert.Nil(err)

	assert.Equal(0, n)
}

func TestSwapCIDForAlias(t *testing.T) {
	assert := assert.New(t)
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

	err = store.UpdateChannel(c)
	assert.Nil(err)

	short, err := store.SwapCIDForAlias(id)
	assert.Nil(err)
	assert.Equal("short1", short)

	_, err = store.SwapCIDForAlias(int64(2))
	assert.NotNil(err)
}

func TestSwapAliasForCID(t *testing.T) {
	assert := assert.New(t)
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

	err = store.UpdateChannel(c)
	assert.Nil(err)

	cid, err := store.SwapAliasForCID(c.Alias)
	assert.Nil(err)
	assert.Equal(int64(1), cid)

	_, err = store.SwapAliasForCID("not_exist")
	assert.NotNil(err)
}

func TestSwapPIDForFilename(t *testing.T) {
	assert := assert.New(t)
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

	podcast, err := store.AddPodcastToChannel(cid, "podcast1.mp3", "podcast1", "10001")
	assert.Nil(err)
	assert.Equal(int64(1), podcast.ID)

	filename, err := store.SwapPIDForFilename(podcast.ID)
	assert.Nil(err)
	assert.Equal("podcast1.mp3", filename)

	_, err = store.SwapPIDForFilename(int64(2))
	assert.NotNil(err)
}
