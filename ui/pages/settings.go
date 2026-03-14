package pages

import (
	"fmt"
	"strconv"

	"holy-codex/infrastructure/config"
	"holy-codex/infrastructure/network"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SettingsPage lets the scribe configure their Codex.
type SettingsPage struct {
	cfg       *config.Config
	discovery *network.Discovery
}

// NewSettingsPage constructs a SettingsPage.
func NewSettingsPage(cfg *config.Config, discovery *network.Discovery) *SettingsPage {
	return &SettingsPage{cfg: cfg, discovery: discovery}
}

// Build returns the Fyne canvas object for this page.
func (p *SettingsPage) Build() fyne.CanvasObject {
	// ── Identity ──────────────────────────────────────────────────────────────
	nameEntry := widget.NewEntry()
	nameEntry.SetText(p.cfg.UserName)

	portEntry := widget.NewEntry()
	portEntry.SetText(strconv.Itoa(p.cfg.NetworkPort))

	syncCheck := widget.NewCheck("Enable LAN Sync", func(on bool) {
		p.cfg.SyncEnabled = on
	})
	syncCheck.SetChecked(p.cfg.SyncEnabled)

	saveBtn := widget.NewButton("✦ Save Edicts", func() {
		p.cfg.UserName = nameEntry.Text
		if port, err := strconv.Atoi(portEntry.Text); err == nil {
			p.cfg.NetworkPort = port
		}
		if err := p.cfg.Save(); err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		dialog.ShowInformation("Saved", "Thy edicts hath been inscribed.", fyne.CurrentApp().Driver().AllWindows()[0])
	})

	// ── Peers ─────────────────────────────────────────────────────────────────
	peersList := widget.NewList(
		func() int { return len(p.discovery.Peers()) },
		func() fyne.CanvasObject {
			return widget.NewLabel("peer")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			peers := p.discovery.Peers()
			if i >= len(peers) {
				return
			}
			peer := peers[i]
			o.(*widget.Label).SetText(fmt.Sprintf("⊕ %s  (%s:%d)", peer.Name, peer.Addr, peer.Port))
		},
	)

	refreshPeers := widget.NewButton("↻ Seek Brethren", func() {
		peersList.Refresh()
	})

	// ── Layout ────────────────────────────────────────────────────────────────
	identity := widget.NewCard("Identity of the Scribe", "",
		container.NewVBox(
			formRow("Name", nameEntry),
			formRow("Sync Port", portEntry),
			syncCheck,
			saveBtn,
		),
	)

	peersCard := widget.NewCard("Known Scribes on the Network", "", container.NewVBox(
		refreshPeers,
		container.NewVScroll(peersList),
	))

	return container.NewVScroll(container.NewVBox(identity, peersCard))
}

func formRow(label string, w fyne.CanvasObject) fyne.CanvasObject {
	return container.NewGridWithColumns(2,
		widget.NewLabelWithStyle(label, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		w,
	)
}