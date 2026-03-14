package pages

import (
	"fmt"

	"holy-codex/infrastructure/storage"
	"holy-codex/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// HomePage shows a welcome scroll and recent entries.
type HomePage struct {
	store    storage.Storage
	notifier *services.Notifier
}

// NewHomePage constructs a HomePage.
func NewHomePage(store storage.Storage, notifier *services.Notifier) *HomePage {
	return &HomePage{store: store, notifier: notifier}
}

// Build returns the Fyne canvas object for this page.
func (p *HomePage) Build() fyne.CanvasObject {
	// ── Header ────────────────────────────────────────────────────────────────
	titleLabel := widget.NewLabelWithStyle(
		" The Holy Codex ",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)

	subtitle := widget.NewLabelWithStyle(
		"Every soul has a story — inscribe yours.",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	separator := widget.NewSeparator()

	// ── Recent entries list ───────────────────────────────────────────────────
	entries, _ := p.store.ListEntries("default")

	list := widget.NewList(
		func() int { return len(entries) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("📜", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabel("Entry Title"),
				widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(entries) {
				return
			}
			e := entries[i]
			row := o.(*fyne.Container)
			row.Objects[1].(*widget.Label).SetText(e.Title)
			row.Objects[2].(*widget.Label).SetText(e.DateLabel())
		},
	)

	// ── Stats strip ───────────────────────────────────────────────────────────
	statsLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("Scrolls written: %d", len(entries)),
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	// ── Layout ────────────────────────────────────────────────────────────────
	header := container.NewVBox(titleLabel, subtitle, separator, statsLabel, widget.NewSeparator())

	return container.NewBorder(header, nil, nil, nil, list)
}