package main

import (
	"log"
	"holy-codex/app"
)

func main() {
	a ,err := app.New()
	if err != nil {
		log.Fatalf("Failed to initilize app: %v", err)
	}
	a.Run()
}