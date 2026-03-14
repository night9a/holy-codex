package pages

import (
	"strings"
	"time"

	"holy-codex/domain"
	"holy-codex/infrastructure/storage"
	"holy-codex/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// DayPage is the writing surface — a parchment scroll the user writes upon.
type DayPage struct {
	store    storage.Storage
	autoSave *services.AutoSave

	titleEntry *widget.Entry
	bodyEntry  *widget.Entry
	tagsEntry  *widget.Entry
	moodSelect *widget.Select
}

// NewDayPage constructs a DayPage.
func NewDayPage(store storage.Storage, autoSave *services.AutoSave) *DayPage {
	return &DayPage{store: store, autoSave: autoSave}
}

// Build returns the Fyne canvas object for this page.
func (p *DayPage) Build() fyne.CanvasObject {
	// ── Widgets ───────────────────────────────────────────────────────────────
	p.titleEntry = widget.NewEntry()
	p.titleEntry.SetPlaceHolder("Give this scroll a title…")

	p.bodyEntry = widget.NewMultiLineEntry()
	p.bodyEntry.SetPlaceHolder(
		"Inscribe thy thoughts here, wanderer.\n\nLet words flow like ink upon ancient vellum…",
	)
	p.bodyEntry.Wrapping = fyne.TextWrapWord

	p.tagsEntry = widget.NewEntry()
	p.tagsEntry.SetPlaceHolder("Tags, separated by commas: travel, dream, gratitude…")

	moodOptions := []string{
		string(domain.MoodJoyful),
		string(domain.MoodCalm),
		string(domain.MoodPensive),
		string(domain.MoodAnxious),
		string(domain.MoodSad),
		string(domain.MoodAngry),
		string(domain.MoodGrateful),
	}
	p.moodSelect = widget.NewSelect(moodOptions, nil)
	p.moodSelect.SetSelected(string(domain.MoodCalm))

	// Auto-save on body change
	p.bodyEntry.OnChanged = func(s string) {
		p.autoSave.Stage(p.buildEntry())
	}

	// ── Date banner ───────────────────────────────────────────────────────────
	dateBanner := widget.NewLabelWithStyle(
		"⊕  "+time.Now().Format("Monday, 2 January 2006"),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)

	// ── Toolbar ───────────────────────────────────────────────────────────────
	saveBtn := widget.NewButton("✦ Seal This Scroll", func() {
		entry := p.buildEntry()
		if err := entry.Validate(); err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		if err := p.store.SaveEntry(entry); err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		p.clearForm()
		dialog.ShowInformation("Sealed", "Thy scroll hath been committed to the Codex.", fyne.CurrentApp().Driver().AllWindows()[0])
	})

	clearBtn := widget.NewButton("✧ Fresh Scroll", func() {
		p.clearForm()
	})

	toolbar := container.NewHBox(saveBtn, clearBtn)

	// ── Form layout ───────────────────────────────────────────────────────────
	form := container.NewVBox(
		dateBanner,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Title", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		p.titleEntry,
		widget.NewLabelWithStyle("Words", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewVScroll(p.bodyEntry),
		widget.NewLabelWithStyle("Mood", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		p.moodSelect,
		widget.NewLabelWithStyle("Tags", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		p.tagsEntry,
		widget.NewSeparator(),
		toolbar,
	)

	return container.NewScroll(form)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (p *DayPage) buildEntry() *domain.DiaryEntry {
	tags := parseTags(p.tagsEntry.Text)
	return domain.NewDiaryEntry(
		"default",
		p.titleEntry.Text,
		p.bodyEntry.Text,
		tags,
		domain.Mood(p.moodSelect.Selected),
	)
}

func (p *DayPage) clearForm() {
	p.titleEntry.SetText("")
	p.bodyEntry.SetText("")
	p.tagsEntry.SetText("")
	p.moodSelect.SetSelected(string(domain.MoodCalm))
}

func parseTags(raw string) []string {
	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	for _, t := range parts {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}