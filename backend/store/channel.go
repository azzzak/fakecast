package store

//
// Add
//

// AddChannel action
func (s *Store) AddChannel() (int64, error) {
	result, err := s.db.Exec("INSERT INTO channels (title) VALUES ('')")
	if err != nil {
		return 0, &Error{Err: err}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, &Error{Err: err}
	}

	return id, nil
}

//
// List
//

// ListChannels action
func (s *Store) ListChannels() ([]Channel, error) {
	rows, err := s.db.Query("SELECT id, alias,title FROM channels")
	if err != nil {
		return nil, &Error{Err: err}
	}
	defer rows.Close()

	var (
		cs []Channel
		c  Channel
	)

	for rows.Next() {
		rows.Scan(&c.ID, &c.Alias, &c.Title)
		cs = append(cs, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, &Error{Err: err}
	}

	return cs, nil
}

//
// Info
//

// ChannelInfo action
func (s *Store) ChannelInfo(cid int64) (*Channel, error) {
	row := s.db.QueryRow("SELECT id, alias, title, description, image, author FROM channels WHERE id=?", cid)
	var c Channel

	err := row.Scan(&c.ID, &c.Alias, &c.Title, &c.Description, &c.Cover, &c.Author)
	if err != nil {
		return nil, &Error{Err: err}
	}

	return &c, nil
}

//
// Update
//

// UpdateChannel action
func (s *Store) UpdateChannel(c *Channel) error {
	_, err := s.db.Exec("UPDATE channels SET alias=?, title=?, description=?, image=?, author=? WHERE id=?", c.Alias, c.Title, c.Description, c.Cover, c.Author, c.ID)
	if err != nil {
		return &Error{Err: err}
	}

	return nil
}

//
// Delete
//

// DeleteChannel action
func (s *Store) DeleteChannel(channel int64) error {
	tx, _ := s.db.Begin()
	_, err := s.db.Exec("DELETE FROM channels WHERE id=?", channel)
	if err != nil {
		return &Error{Err: err}
	}

	_, err = s.db.Exec("DELETE FROM podcasts WHERE channel=?", channel)
	if err != nil {
		return &Error{Err: err}
	}
	tx.Commit()

	return nil
}
