package app

import (
	"holy-codex/infrastructure/config"
	"holy-codex/infrastructure/network"
	"holy-codex/infrastructure/storage"
	"holy-codex/services"
	"holy-codex/ui/pages"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// App is the root application struct, wiring together all subsystems.
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	ctx     *Context
}

// New constructs and wires the full application.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	store, err := storage.NewSqliteStorage(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	discovery := network.NewDiscovery(cfg.NetworkPort)
	syncer := network.NewSync(store, discovery, cfg.NetworkPort)
	autoSave := services.NewAutoSave(store)
	syncSvc := services.NewSyncService(syncer)
	notifier := services.NewNotifier()

	ctx := &Context{
		Config:      cfg,
		Storage:     store,
		Discovery:   discovery,
		Syncer:      syncer,
		AutoSave:    autoSave,
		SyncService: syncSvc,
		Notifier:    notifier,
	}

	fyneA := fyneapp.New()
	fyneA.Settings().SetTheme(&CodicTheme{})

	win := fyneA.NewWindow("Holy Codex — My Diary")
	win.Resize(fyne.NewSize(960, 680))
	win.CenterOnScreen()

	return &App{
		fyneApp: fyneA,
		window:  win,
		ctx:     ctx,
	}, nil
}

// Run starts background services and opens the main window.
func (a *App) Run() {
	a.ctx.AutoSave.Start()
	a.ctx.SyncService.Start()

	a.window.SetContent(a.buildMain())
	a.window.ShowAndRun()

	// Cleanup on exit
	a.ctx.AutoSave.Stop()
	a.ctx.SyncService.Stop()
	_ = a.ctx.Storage.Close()
}

// buildMain assembles the top-level navigation layout.
func (a *App) buildMain() fyne.CanvasObject {
	homePage := pages.NewHomePage(a.ctx.Storage, a.ctx.Notifier)
	dayPage := pages.NewDayPage(a.ctx.Storage, a.ctx.AutoSave)
	settingsPage := pages.NewSettingsPage(a.ctx.Config, a.ctx.Discovery)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Codex", theme.HomeIcon(), homePage.Build()),
		container.NewTabItemWithIcon("Write", theme.DocumentCreateIcon(), dayPage.Build()),
		container.NewTabItemWithIcon("Scrolls", theme.SettingsIcon(), settingsPage.Build()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	header := widget.NewLabelWithStyle(
		"✦ Holy Codex",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true},
	)

	return container.NewBorder(header, nil, nil, nil, tabs)
}