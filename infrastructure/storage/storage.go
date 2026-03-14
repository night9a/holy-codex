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

	DeleteEntry(id string) error
	UnsyncedEntries() ([]*domain.DiaryEntry,error)
	MarkSynced(id string) error
	Close() error
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

func (s *SQLiteStorage) migrate() error {
	return RunMigrations(s.db)
}

func (s *SQLiteStorage) SaveEntry(e *domain.DiaryEntry) error {
	tags, _ := json.Marshal(e.Tags)
	_, err := s.db.Exec(`
		INSERT INTO entries (id, user_id, title, body, tags, mood, created_at, updated_at, is_synced)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title, body=excluded.body, tags=excluded.tags,
			mood=excluded.mood, updated_at=excluded.updated_at, is_synced=excluded.is_synced
	`, e.ID, e.UserID, e.Title, e.Body, string(tags), string(e.Mood),
		e.CreatedAt.UTC(), e.UpdatedAt.UTC(), e.IsSynced)
	return err
}

func (s *SQLiteStorage) GetEntry(id string) (*domain.DiaryEntry, error) {
	row := s.db.QueryRow(`SELECT id, user_id, title, body, tags, mood, created_at, updated_at, is_synced FROM entries WHERE id=?`, id)
	return scanEntry(row)
}

func (s *SQLiteStorage) ListEntries(userID string) ([]*domain.DiaryEntry, error) {
	rows, err := s.db.Query(`SELECT id, user_id, title, body, tags, mood, created_at, updated_at, is_synced FROM entries WHERE user_id=? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	var entries []*domain.DiaryEntry
	for rows.Next() {
		e, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
 
func (s *SQLiteStorage) DeleteEntry(id string) error {
	_, err := s.db.Exec(`DELETE FROM entries WHERE id=?`, id)
	return err
}
 
func (s *SQLiteStorage) UnsyncedEntries() ([]*domain.DiaryEntry, error) {
	rows, err := s.db.Query(`SELECT id, user_id, title, body, tags, mood, created_at, updated_at, is_synced FROM entries WHERE is_synced=0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	var entries []*domain.DiaryEntry
	for rows.Next() {
		e, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
 
func (s *SQLiteStorage) MarkSynced(id string) error {
	_, err := s.db.Exec(`UPDATE entries SET is_synced=1 WHERE id=?`, id)
	return err
}
 
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
 
// ─── Scan helper ──────────────────────────────────────────────────────────────
 
type scanner interface {
	Scan(dest ...any) error
}
 
func scanEntry(row scanner) (*domain.DiaryEntry, error) {
	var e domain.DiaryEntry
	var tagsJSON string
	var createdAt, updatedAt time.Time
	var mood string
 
	err := row.Scan(&e.ID, &e.UserID, &e.Title, &e.Body, &tagsJSON, &mood, &createdAt, &updatedAt, &e.IsSynced)
	if err != nil {
		return nil, err
	}
	e.Mood = domain.Mood(mood)
	e.CreatedAt = createdAt
	e.UpdatedAt = updatedAt
	_ = json.Unmarshal([]byte(tagsJSON), &e.Tags)
	return &e, nil
}