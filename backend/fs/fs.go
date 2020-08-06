package fs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CoverDirName const
const CoverDirName = "cover"

// PodcastsDirName const
const PodcastsDirName = "podcasts"

// FrontDirName const
const FrontDirName = "front"

// Error type
type Error struct {
	Err error
}

func (err *Error) Error() string {
	return err.Err.Error()
}

// Dir entity
type Dir struct {
	Root string
}

// NewRoot constructor
func NewRoot(root string) Dir {
	return Dir{
		Root: filepath.Join(root, PodcastsDirName),
	}
}

//
// Create
//

// CreateDir action
func (d *Dir) CreateDir(channel int64) error {
	path := filepath.Join(d.Root, strconv.FormatInt(channel, 10), CoverDirName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return &Error{Err: err}
	}
	return nil
}

//
// Remove
//

// RemoveDir action
func (d *Dir) RemoveDir(channel string) error {
	path := filepath.Join(d.Root, channel)
	if err := os.RemoveAll(path); err != nil {
		return &Error{Err: err}
	}
	return nil
}

// RemovePodcast action
func (d *Dir) RemovePodcast(channel, podcast string) error {
	path := filepath.Join(d.Root, channel, podcast)
	return d.remove(path)
}

// RemoveCover action
func (d *Dir) RemoveCover(channel, cover string) error {
	path := filepath.Join(d.Root, channel, CoverDirName, cover)
	return d.remove(path)
}

func (d *Dir) remove(path string) error {
	if err := os.Remove(path); err != nil {
		return &Error{Err: err}
	}
	return nil
}

//
// Save
//

// SavePodcastToDir action
func (d *Dir) SavePodcastToDir(channel, filename string) (*os.File, error) {
	path := filepath.Join(d.Root, channel, filename)
	return d.save(path)
}

// SaveCover action
func (d *Dir) SaveCover(channel, filename string) (*os.File, error) {
	path := filepath.Join(d.Root, channel, CoverDirName, filename)
	return d.save(path)
}

func (d *Dir) save(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, &Error{Err: err}
	}
	return f, nil
}

//
// Rename
//

// RenameDir action
func (d *Dir) RenameDir(old, new string) error {
	expose := func(s string) string {
		return filepath.Join(d.Root, s)
	}
	if err := os.Rename(expose(old), expose(new)); err != nil {
		return &Error{Err: err}
	}
	return nil
}

//
// Helpers
//

// IsDirExist helper
func (d *Dir) IsDirExist(channel string) bool {
	return isFileExist(d.Root, channel)
}

// IsPodcastExist helper
func (d *Dir) IsPodcastExist(channel, filename string) bool {
	return isFileExist(d.Root, channel, filename)
}

func isFileExist(p ...string) bool {
	path := filepath.Join(p...)
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// NameAndExtFrom helper
func NameAndExtFrom(filename string) (string, string) {
	t := strings.Split(filename, ".")
	if len(t) > 1 {
		name, ext := t[:len(t)-1], t[len(t)-1]
		return strings.Join(name, "."), ext
	}
	return filename, ""
}
