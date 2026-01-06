package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS transcripts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		url TEXT,
		backend TEXT,
		text TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func (s *Store) Save(t Transcript) error {
	_, err := s.db.Exec(`
	INSERT INTO transcripts (title, url, backend, text)
	VALUES (?, ?, ?, ?)`,
		t.Title, t.URL, t.Backend, t.Text,
	)
	return err
}

func (s *Store) List() ([]Transcript, error) {
	rows, err := s.db.Query(`
	SELECT id, title, url, backend, text, created_at
	FROM transcripts
	ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Transcript
	for rows.Next() {
		var t Transcript
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.URL,
			&t.Backend,
			&t.Text,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM transcripts WHERE id = ?`, id)
	return err
}
