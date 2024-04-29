package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %v", operation, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS urls(
	    id INTEGER PRIMARY KEY, 
	    alias TEXT NOT NULL UNIQUE, 
	    url TEXT NOT NULL UNIQUE);
	CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to create table: %v", operation, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute statement: %v", operation, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) error {
	const operation = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO urls(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %v", operation, err)
	}
	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		// TODO: refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return storage.ErrAliasExists
		}
		return fmt.Errorf("%s: failed to execute statement: %v", operation, err)
	}

	//id, err = res.LastInsertId()
	//if err != nil {
	//	return fmt.Errorf("%s: failed to get last insert id: %v", operation, err)
	//}
	return nil

}

func (s *Storage) GetUrl(alias string) (string, error) {
	const operation = "storage.sqlite.GetUrl"
	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: failed to prepare statement: %v", operation, err)
	}

	var urlReceived string
	err = stmt.QueryRow(alias).Scan(&urlReceived)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrAliasNotFound
		}
		return "", fmt.Errorf("%s: failed to execute statement: %v", operation, err)
	}

	return urlReceived, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const operation = "storage.sqlite.DeleteUrl"
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %v", operation, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: failed to execute statement: %v", operation, err)
	}

	deletedRows, err := res.RowsAffected()
	if err != nil {
		if deletedRows == 0 {
			return storage.ErrAliasNotFound
		}
		return fmt.Errorf("%s: failed to get rows affected: %v", operation, err)
	}

	return nil
}
