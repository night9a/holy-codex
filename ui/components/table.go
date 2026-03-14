package components

import (
	"holy-codex/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// EntryTable renders a scrollable table of diary entries.
type EntryTable struct {
	entries  []*domain.DiaryEntry
	onSelect func(*domain.DiaryEntry)
	table    *widget.Table
}

// NewEntryTable constructs an EntryTable.
// onSelect is called whenever the user taps a row.
func NewEntryTable(entries []*domain.DiaryEntry, onSelect func(*domain.DiaryEntry)) *EntryTable {
	et := &EntryTable{
		entries:  entries,
		onSelect: onSelect,
	}
	et.table = et.build()
	return et
}

// Refresh replaces the dataset and repaints the table.
func (et *EntryTable) Refresh(entries []*domain.DiaryEntry) {
	et.entries = entries
	et.table.Refresh()
}

// CanvasObject returns the underlying Fyne widget.
func (et *EntryTable) CanvasObject() fyne.CanvasObject {
	return container.NewVScroll(et.table)
}

// ─── Internal ─────────────────────────────────────────────────────────────────

var tableHeaders = []string{"Date", "Title", "Mood", "Tags"}

func (et *EntryTable) build() *widget.Table {
	t := widget.NewTable(
		// dimensions: rows = header + entries, cols = 4
		func() (int, int) { return len(et.entries) + 1, len(tableHeaders) },

		// create cell template
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
		},

		// update cell content
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 {
				// Header row
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.SetText(tableHeaders[id.Col])
				return
			}
			e := et.entries[id.Row-1]
			label.TextStyle = fyne.TextStyle{}
			switch id.Col {
			case 0:
				label.SetText(e.DateLabel())
			case 1:
				label.SetText(truncate(e.Title, 40))
			case 2:
				label.SetText(moodGlyph(e.Mood) + " " + string(e.Mood))
			case 3:
				label.SetText(joinTags(e.Tags, 3))
			}
		},
	)

	// Column widths
	t.SetColumnWidth(0, 130)
	t.SetColumnWidth(1, 280)
	t.SetColumnWidth(2, 110)
	t.SetColumnWidth(3, 200)

	t.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 || et.onSelect == nil {
			return
		}
		et.onSelect(et.entries[id.Row-1])
	}

	return t
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func joinTags(tags []string, max int) string {
	out := ""
	for i, t := range tags {
		if i >= max {
			out += " …"
			break
		}
		if i > 0 {
			out += ", "
		}
		out += t
	}
	return out
}

func moodGlyph(m domain.Mood) string {
	switch m {
	case domain.MoodJoyful:
		return "☀"
	case domain.MoodCalm:
		return "🌿"
	case domain.MoodPensive:
		return "🌙"
	case domain.MoodAnxious:
		return "⚡"
	case domain.MoodSad:
		return "🌧"
	case domain.MoodAngry:
		return "🔥"
	case domain.MoodGrateful:
		return "✦"
	}
	return "·"
}