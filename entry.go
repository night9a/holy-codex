package main

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

type Entry struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Content   string    `json:"content"`
}

func NewEntry(content string) Entry {
    return Entry{
        ID:        uuid.New().String(),
        Timestamp: time.Now(),
        Content:   content,
    }
}

func (e Entry) ToJSON() ([]byte, error) {
    return json.Marshal(e)
}

func EntryFromJSON(data []byte) (Entry, error) {
    var e Entry
    err := json.Unmarshal(data, &e)
    return e, err
}
