package store

import (
	"strconv"
)

const (
	shortForm = iota
	fullForm
)

//
// Add
//

// AddPodcastToChannel action
func (s *Store) AddPodcastToChannel(cid int64, filename, title, size string) (*Podcast, error) {
	length, err := strconv.Atoi(size)
	if err != nil {
		return nil, &Error{Err: err}
	}

	result, err := s.db.Exec("INSERT INTO podcasts (channel, filename, title, length) VALUES (?, ?, ?, ?)", cid, filename, title, length)
	if err != nil {
		return nil, &Error{Err: err}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, &Error{Err: err}
	}

	p := Podcast{
		ID:    id,
		Title: title,
	}

	return &p, nil
}

//
// Info
//

// PodcastInfo action
func (s *Store) PodcastInfo(pid int64) (*Podcast, error) {
	row := s.db.QueryRow("SELECT id, filename, published, title, length, guid, pub_date, description, duration, image, explicit, season, episode FROM podcasts WHERE id=?", pid)
	var p Podcast

	err := row.Scan(&p.ID, &p.Filename, &p.Published, &p.Title, &p.Length, &p.GUID, &p.PubDate, &p.Description, &p.Duration, &p.Artwork, &p.Explicit, &p.Season, &p.Episode)
	if err != nil {
		return nil, &Error{Err: err}
	}

	return &p, nil
}

//
// List
//

// ListPodcastsFrom action
func (s *Store) ListPodcastsFrom(cid int64) ([]Podcast, error) {
	return s.listPodcasts(cid, shortForm)
}

// ListFullPodcastsFrom action
func (s *Store) ListFullPodcastsFrom(cid int64) ([]Podcast, error) {
	return s.listPodcasts(cid, fullForm)
}

// ListPodcastsFrom action
func (s *Store) listPodcasts(cid int64, form int) ([]Podcast, error) {
	var (
		sql     string
		holders []interface{}
	)

	var (
		podcasts []Podcast
		p        Podcast
	)

	switch form {
	case shortForm:
		sql = "SELECT id, filename, title FROM podcasts WHERE channel=? ORDER BY id DESC"
		holders = []interface{}{&p.ID, &p.Filename, &p.Title}
	case fullForm:
		sql = "SELECT id, filename, published, title, length, guid, pub_date, description, duration, image, explicit, season, episode FROM podcasts WHERE channel=? ORDER BY id DESC"
		holders = []interface{}{&p.ID, &p.Filename, &p.Published, &p.Title, &p.Length, &p.GUID, &p.PubDate, &p.Description, &p.Duration, &p.Artwork, &p.Explicit, &p.Season, &p.Episode}
	}

	rows, err := s.db.Query(sql, cid)
	if err != nil {
		return nil, &Error{Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(holders...)
		podcasts = append(podcasts, p)
	}

	err = rows.Err()
	if err != nil {
		return nil, &Error{Err: err}
	}

	return podcasts, nil
}

//
// Update
//

// UpdatePodcast action
func (s *Store) UpdatePodcast(p *Podcast) error {
	_, err := s.db.Exec("UPDATE podcasts SET published=?, title=?, length=?, guid=?, pub_date=?, description=?, duration=?, image=?, explicit=?, season=?, episode=? WHERE id=?", p.Published, p.Title, p.Length, p.GUID, p.PubDate, p.Description, p.Duration, p.Artwork, p.Explicit, p.Season, p.Episode, p.ID)
	if err != nil {
		return &Error{Err: err}
	}

	return nil
}

//
// Delete
//

// DeletePodcast action
func (s *Store) DeletePodcast(pid int64) error {
	_, err := s.db.Exec("DELETE FROM podcasts WHERE id=?", pid)
	if err != nil {
		return &Error{Err: err}
	}

	return nil
}
