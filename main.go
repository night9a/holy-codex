package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func main() {
    a := app.New()
    w := a.NewWindow("Holy Codex")

    // واجهة بسيطة: نص وزر
    label := widget.NewLabel("Welcome to Holy Codex")
    button := widget.NewButton("Click me", func() {
        label.SetText("Button clicked!")
    })

    w.SetContent(container.NewVBox(label, button))
    w.ShowAndRun()
}
