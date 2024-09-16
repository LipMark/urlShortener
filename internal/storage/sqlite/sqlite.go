package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"urlShortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// CreateStorage create a new storage for data entries
func CreateStorage(storagePath string) (*Storage, error) {
	const errstring = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errstring, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		urlalias TEXT NOT NULL UNIQUE,
		originalurl TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errstring, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errstring, err)
	}

	return &Storage{db: db}, nil
}

// SaveURL - add an entry into DB
func (s *Storage) SaveURL(originalURL string, alias string) (int64, error) {
	const errstringSaveURL = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(originalurl, urlalias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errstringSaveURL, err)
	}

	res, err := stmt.Exec(originalURL, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", errstringSaveURL, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", errstringSaveURL, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", errstringSaveURL, err)
	}

	return id, nil
}

// GetURL - get an entry from DB
func (s *Storage) GetURL(alias string) (string, error) {
	const errstringGetURL = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT originalurl FROM url WHERE urlalias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", errstringGetURL, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", errstringGetURL, err)
	}

	return resURL, nil
}

// DeleteURL - delete an entry in DB
func (s *Storage) DeleteURL(alias string) error {
	const errstringDeleteURL = "storage.sqlite.DeleteURL"
	stmt, err := s.db.Prepare("DELETE FROM url WHERE urlalias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", errstringDeleteURL, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", errstringDeleteURL, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", errstringDeleteURL, err)
	}

	if n == 0 {
		return fmt.Errorf("%s: %w", errstringDeleteURL, storage.ErrURLNotFound)
	}

	return nil
}
