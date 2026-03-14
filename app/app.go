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

type App struct {
	fyneApp fyne.App
	window fyne.Window
	ctx *Context
}

func New() (*App,error) {
	cfg, err := config.Load()
	if err != nil {
		return nil,err
	}
	store, err := storage.NewSqliteStorage(cfg.DBPath)
	if err != nil {
		return nil,err
	}
	discovery := network.NewDiscovery(cfg.NetworkPort)
	syncer := network.NewSync(store,discovery,cfg.NetworkPort)
	autoSave := services.NewAutoSave(store)
	syncSvc := services.NewSyncService(syncer)
	notifier := services.NewNotifier()

	ctx := &Context{
		Config: cfg,
		Storage: store,
		Discovery: discovery,
		Syncer: syncer,
		autoSave: autoSave,
		SyncService: syncSvc,
		Notifier: notifier,
	}

	fyneA := fyneapp.new()
	fyneA.Settings().SetTheme(&CodicTheme{})

	win := fyneA.NewWindow("Holy Codex")
	win.Resize(fyne.NewSize(960,680))
	win.CenterOnScreen()

	return &App{
		fyneApp: fyneA,
		window: win,
		ctx: ctx,
	},nil
}
func (a *App) buildMain() fyne.CanvasObject {
	HomePage := pages.NewHomePage(a.ctx.Storage,a.ctx.Notifier)

	//add other pages

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Codex",theme.HomeIcon(),HomePage.Build()),
		//add otherpages
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	header := widget.NewLabelWithStyle(
		"# Holy Codex",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true,Italic: true},
	)

	return container.NewBorder(header,nil,nil,nil.tabs)
}