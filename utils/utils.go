package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// ─── Date helpers ─────────────────────────────────────────────────────────────

// ParseDate attempts to parse a date string using several common layouts.
func ParseDate(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"2 January 2006",
		"January 2, 2006",
		time.RFC3339,
	}
	s = strings.TrimSpace(s)
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date: %q", s)
}

// FormatDate returns a human-friendly date string.
func FormatDate(t time.Time) string {
	return t.Format("Monday, 2 January 2006")
}

// StartOfDay returns midnight of the given day in local time.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns 23:59:59 of the given day in local time.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// ─── Logger ───────────────────────────────────────────────────────────────────

// Logger wraps the standard library logger with level prefixes.
type Logger struct {
	inner *log.Logger
	level Level
}

// Level is the minimum log severity to emit.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// NewLogger creates a Logger writing to stdout at the given minimum level.
func NewLogger(level Level) *Logger {
	return &Logger{
		inner: log.New(os.Stdout, "", log.LstdFlags),
		level: level,
	}
}

func (l *Logger) Debug(format string, args ...any) {
	if l.level <= LevelDebug {
		l.inner.Printf("[DEBUG] "+format, args...)
	}
}

func (l *Logger) Info(format string, args ...any) {
	if l.level <= LevelInfo {
		l.inner.Printf("[INFO]  "+format, args...)
	}
}

func (l *Logger) Warn(format string, args ...any) {
	if l.level <= LevelWarn {
		l.inner.Printf("[WARN]  "+format, args...)
	}
}

func (l *Logger) Error(format string, args ...any) {
	if l.level <= LevelError {
		l.inner.Printf("[ERROR] "+format, args...)
	}
}

// ─── String helpers ───────────────────────────────────────────────────────────

// Truncate clips s to maxLen runes, appending "…" if clipped.
func Truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}

// ContainsAny returns true if s contains any of the substrings (case-insensitive).
func ContainsAny(s string, subs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range subs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}