package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "time"
)

var peers []string // قائمة الأجهزة الأخرى على LAN

// اكتشاف الأجهزة عن طريق UDP broadcast
func BroadcastPresence(myID string) {
    addr := net.UDPAddr{
        IP:   net.IPv4bcast,
        Port: 9999,
    }
    conn, _ := net.DialUDP("udp4", nil, &addr)
    defer conn.Close()

    for {
        conn.Write([]byte("HolyCodexNode:" + myID))
        time.Sleep(5 * time.Second)
    }
}

// إرسال entry لكل peer
func SyncEntry(entry Entry) {
    jsonData, _ := entry.ToJSON()
    for _, peer := range peers {
        go http.Post(fmt.Sprintf("http://%s:8080/sync", peer), "application/json", bytes.NewBuffer(jsonData))
    }
}

// استقبال entry من peer
func HandleSync(w http.ResponseWriter, r *http.Request) {
    var e Entry
    json.NewDecoder(r.Body).Decode(&e)
    SaveEntry(e)
    AddToArchive(e)
    w.WriteHeader(200)
}
