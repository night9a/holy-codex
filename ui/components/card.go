package components

import (
	"holy-codex/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// EntryCard renders a single diary entry as a parchment-style card.
type EntryCard struct {
	entry    *domain.DiaryEntry
	onEdit   func(*domain.DiaryEntry)
	onDelete func(*domain.DiaryEntry)
}

// NewEntryCard constructs a card for the given entry.
func NewEntryCard(
	entry *domain.DiaryEntry,
	onEdit func(*domain.DiaryEntry),
	onDelete func(*domain.DiaryEntry),
) *EntryCard {
	return &EntryCard{entry: entry, onEdit: onEdit, onDelete: onDelete}
}

// CanvasObject builds and returns the card UI.
func (c *EntryCard) CanvasObject() fyne.CanvasObject {
	e := c.entry

	// ── Title & date ──────────────────────────────────────────────────────────
	title := widget.NewLabelWithStyle(
		e.Title,
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	date := widget.NewLabelWithStyle(
		e.DateLabel()+"  "+moodGlyph(e.Mood)+" "+string(e.Mood),
		fyne.TextAlignLeading,
		fyne.TextStyle{Italic: true},
	)

	// ── Body preview ──────────────────────────────────────────────────────────
	preview := widget.NewLabel(truncate(e.Body, 160))
	preview.Wrapping = fyne.TextWrapWord

	// ── Tags ──────────────────────────────────────────────────────────────────
	tagStr := joinTags(e.Tags, 5)
	tagsLabel := widget.NewLabelWithStyle(
		"🏷 "+tagStr,
		fyne.TextAlignLeading,
		fyne.TextStyle{Italic: true},
	)

	// ── Sync indicator ────────────────────────────────────────────────────────
	syncBadge := "⊙ local"
	if e.IsSynced {
		syncBadge = "✦ synced"
	}
	syncLabel := widget.NewLabelWithStyle(syncBadge, fyne.TextAlignTrailing, fyne.TextStyle{})

	// ── Actions ───────────────────────────────────────────────────────────────
	editBtn := widget.NewButton("✎ Edit", func() {
		if c.onEdit != nil {
			c.onEdit(e)
		}
	})
	deleteBtn := widget.NewButton("✗ Burn", func() {
		if c.onDelete != nil {
			c.onDelete(e)
		}
	})

	actions := container.NewHBox(editBtn, deleteBtn, widget.NewLabel(""), syncLabel)

	inner := container.NewVBox(
		title,
		date,
		widget.NewSeparator(),
		preview,
		tagsLabel,
		widget.NewSeparator(),
		actions,
	)

	card := widget.NewCard("", "", inner)
	return card
}