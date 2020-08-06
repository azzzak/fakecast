package store

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // sqlite
)

const storeFile = "fakecast.db"

// Store entity
type Store struct {
	db *sql.DB
}

// Channel entity
type Channel struct {
	ID          int64  `json:"id"`
	Alias       string `json:"alias"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Cover       string `json:"cover"`
	Author      string `json:"author,omitempty"`
	Host        string `json:"host"`
}

// Podcast entity
type Podcast struct {
	ID          int64  `json:"id"`
	Filename    string `json:"filename"`
	Published   int    `json:"published"`
	Title       string `json:"title"`
	Length      int    `json:"length"`
	GUID        string `json:"guid,omitempty"`
	PubDate     string `json:"pub_date,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	Artwork     string `json:"artwork"`
	Explicit    int    `json:"explicit"`
	Season      int    `json:"season,omitempty"`
	Episode     int    `json:"episode,omitempty"`
}

// Error type
type Error struct {
	Err error
}

func (err *Error) Error() string {
	return err.Err.Error()
}

// NewStore constructor
func NewStore(root string) (Store, error) {
	database, err := sql.Open("sqlite3", filepath.Join(root, storeFile))
	if err != nil {
		return Store{}, &Error{Err: err}
	}

	tx, err := database.Begin()
	if err != nil {
		return Store{}, &Error{Err: err}
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id INTEGER PRIMARY KEY,
			alias TEXT UNIQUE DEFAULT '',
			title TEXT,
			description TEXT DEFAULT '',
			image TEXT DEFAULT '',
			explicit INTEGER DEFAULT 0,
			author TEXT DEFAULT ''
		)
		`)
	if err != nil {
		return Store{}, &Error{Err: err}
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS podcasts (
			id INTEGER PRIMARY KEY, 
			channel INTEGER, 
			filename TEXT, 
			published INTEGER DEFAULT 0,
			title TEXT, 
			length INTEGER,
			guid TEXT DEFAULT '',
			pub_date TEXT DEFAULT '', 
			description TEXT DEFAULT '', 
			duration INTEGER DEFAULT 0,
			image TEXT DEFAULT '',
			explicit INTEGER DEFAULT 0,
			season INTEGER DEFAULT 0, 
			episode INTEGER DEFAULT 0
		)
		`)
	err = tx.Commit()
	if err != nil {
		return Store{}, &Error{Err: err}
	}

	return Store{db: database}, nil
}

// SwapCIDForAlias exchange CID for alias
func (s *Store) SwapCIDForAlias(cid int64) (short string, err error) {
	res, err := s.swap("SELECT alias FROM channels WHERE id=?", cid, short)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

// SwapAliasForCID exchange alias for CID
func (s *Store) SwapAliasForCID(alias string) (cid int64, err error) {
	res, err := s.swap("SELECT id FROM channels WHERE alias=?", alias, cid)
	if err != nil {
		return 0, err
	}
	return res.(int64), err
}

// SwapPIDForFilename PID for filename
func (s *Store) SwapPIDForFilename(pid int64) (filename string, err error) {
	res, err := s.swap("SELECT filename FROM podcasts WHERE id=?", pid, filename)
	if err != nil {
		return "", err
	}
	return res.(string), err
}

func (s *Store) swap(sql string, p, holder interface{}) (interface{}, error) {
	row := s.db.QueryRow(sql, p)
	err := row.Scan(&holder)
	if err != nil {
		return nil, &Error{Err: err}
	}
	return holder, nil
}

// DropChannels table
func (s *Store) DropChannels() error {
	_, err := s.db.Exec("DROP TABLE IF EXISTS channels")
	if err != nil {
		return &Error{Err: err}
	}
	return nil
}

// DropPodcasts table
func (s *Store) DropPodcasts() error {
	_, err := s.db.Exec("DROP TABLE IF EXISTS podcasts")
	if err != nil {
		return &Error{Err: err}
	}
	return nil
}

// Close DB connection
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return &Error{Err: err}
	}
	return nil
}
