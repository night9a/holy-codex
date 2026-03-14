package components

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Calendar is a simple monthly calendar widget that highlights days
// that have diary entries.
type Calendar struct {
	year     int
	month    time.Month
	marked   map[int]bool // day number -> has entry
	onSelect func(time.Time)

	container fyne.CanvasObject
}

// NewCalendar builds a Calendar for the given month.
// markedDays is the set of day-numbers that should be highlighted.
// onSelect is called when a day cell is tapped.
func NewCalendar(year int, month time.Month, markedDays []int, onSelect func(time.Time)) *Calendar {
	c := &Calendar{
		year:     year,
		month:    month,
		marked:   make(map[int]bool),
		onSelect: onSelect,
	}
	for _, d := range markedDays {
		c.marked[d] = true
	}
	c.container = c.build()
	return c
}

// CanvasObject returns the widget tree.
func (c *Calendar) CanvasObject() fyne.CanvasObject {
	return c.container
}

// ─── Internal ─────────────────────────────────────────────────────────────────

var weekdays = []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}

func (c *Calendar) build() fyne.CanvasObject {
	// Month/year header
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("%s  %d", c.month.String(), c.year),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)

	// Weekday row
	dayRow := make([]fyne.CanvasObject, 7)
	for i, d := range weekdays {
		dayRow[i] = widget.NewLabelWithStyle(d, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}
	weekHeader := container.NewGridWithColumns(7, dayRow...)

	// Day grid
	first := time.Date(c.year, c.month, 1, 0, 0, 0, 0, time.Local)
	offset := int(first.Weekday()+6) % 7 // Monday = 0
	daysInMonth := daysIn(c.year, c.month)

	cells := make([]fyne.CanvasObject, 0, offset+daysInMonth)

	// Empty leading cells
	for i := 0; i < offset; i++ {
		cells = append(cells, widget.NewLabel(""))
	}

	// Day buttons
	for day := 1; day <= daysInMonth; day++ {
		d := day // capture
		label := fmt.Sprintf("%d", d)
		if c.marked[d] {
			label = "✦" + label
		}
		btn := widget.NewButton(label, func() {
			if c.onSelect != nil {
				c.onSelect(time.Date(c.year, c.month, d, 0, 0, 0, 0, time.Local))
			}
		})
		cells = append(cells, btn)
	}

	grid := container.NewGridWithColumns(7, cells...)

	return container.NewVBox(header, widget.NewSeparator(), weekHeader, grid)
}

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}