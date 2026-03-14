package domain

import (
	"errors"
	"time"

	"github/google/uuid"
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
	ID	string ''
}