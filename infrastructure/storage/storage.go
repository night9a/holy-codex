package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	"holy-codex/domain"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	SaveEntry(entry *domain.DiaryEntry) error
	GetEntry(id string) (*domain.DiaryEntry,error)
	ListEntries(userID string) ([]*domain.DiaryEntry,error)
}


type SQLiteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(path string) (*SQLiteStorage,error) {
	db,err := sql.Open("sqlite3",path)
	if err != nil {
		return nil,err
	}
	s := &SQLiteStorage{db: db}
	if err := s.migrate(); err != nil {
		return nil,err
	}
	return s,nil
}

