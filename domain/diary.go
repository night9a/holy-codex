package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Mood string

const (
	MoodJoyful Mood = "joyful"
	MoodCalm Mood = "calm"
	MoodPensive Mood = "pensive"
	MoodAnxious Mood = "anxious"
	MoodSad Mood = "sad"
	MoodAngry Mood = "angry"
	MoodGrateful Mood = "grateful"
)

type DiaryEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	Mood      Mood      `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsSynced  bool      `json:"is_synced"`
}

func NewDiaryEntry(userID,title,body string,tags []string,mood Mood) *DiaryEntry {
	now := time.Now()
	return &DiaryEntry{
		ID:	uuid.NewString(),
		UserID: userID,
		Title: title,
		Body: body,
		Tags: tags,
		Mood: mood,
		CreatedAt: now,
		UpdatedAt: now,
		IsSynced: false,
	}
}

func (e *DiaryEntry) Validate() error {
	if e.ID == "" {
		return errors.New("Diary entry must have an ID")
	}
	if len(e.Body) == 0 {
		return errors.New("diary entry body cannot be empty")
	}
	if e.CreatedAt.IsZero() {
		return errors.New("diary entry must have a creation time")
	}
	return nil
}

func (e *DiaryEntry) Touch() {
	e.UpdatedAt = time.Now()
	e.IsSynced = false
}

func (e *DiaryEntry) DateLabel() string {
	return e.CreatedAt.Format("2 january 2006")
}