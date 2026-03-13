package main

import (
    "fmt"
    "net/http"
    "time"

    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

var Archive = make(map[int]map[time.Month]map[int][]Entry)
var myDeviceID = fmt.Sprintf("Device-%d", time.Now().UnixNano())

func AddToArchive(e Entry) {
    y, m, _ := e.Timestamp.Date()
    _, w := e.Timestamp.ISOWeek()

    if _, ok := Archive[y]; !ok {
        Archive[y] = make(map[time.Month]map[int][]Entry)
    }
    if _, ok := Archive[y][m]; !ok {
        Archive[y][m] = make(map[int][]Entry)
    }

    Archive[y][m][w] = append(Archive[y][m][w], e)
}

func main() {
    // تحميل اليوميات الموجودة مسبقًا
    entries, _ := LoadEntries()
    for _, e := range entries {
        AddToArchive(e)
    }

    // بدء GUI
    a := app.New()
    w := a.NewWindow("Holy Codex")

    entryBox := widget.NewMultiLineEntry()
    saveBtn := widget.NewButton("Save", func() {
        e := NewEntry(entryBox.Text)
        SaveEntry(e)
        AddToArchive(e)
        SyncEntry(e) // أرسل لكل peer
        entryBox.SetText("")
    })

    w.SetContent(container.NewVBox(entryBox, saveBtn))

    // بدء HTTP server لاستقبال updates من peers
    go func() {
        http.HandleFunc("/sync", HandleSync)
        http.ListenAndServe(":8080", nil)
    }()

    // بدء broadcast للإعلان عن الجهاز
    go BroadcastPresence(myDeviceID)

    w.ShowAndRun()
}
