package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("hello")

	w.SetContent(widget.NewButton(
		"click me",
		func() { println("hello!")},
	))
	w.ShowAndRun()
}