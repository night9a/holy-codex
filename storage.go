package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
)

const DiaryFile = "diary.json"

// حفظ entry جديد
func SaveEntry(e Entry) error {
    var entries []Entry

    if _, err := os.Stat(DiaryFile); err == nil {
        data, _ := ioutil.ReadFile(DiaryFile)
        json.Unmarshal(data, &entries)
    }

    entries = append(entries, e)
    data, err := json.MarshalIndent(entries, "", "  ")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(DiaryFile, data, 0644)
}

// قراءة كل اليوميات
func LoadEntries() ([]Entry, error) {
    var entries []Entry
    if _, err := os.Stat(DiaryFile); err != nil {
        return entries, nil
    }

    data, err := ioutil.ReadFile(DiaryFile)
    if err != nil {
        return nil, err
    }

    json.Unmarshal(data, &entries)
    return entries, nil
}
