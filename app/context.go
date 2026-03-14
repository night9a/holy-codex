package app

import (
	"holy-codex/infrastructure/config"
	"holy-codex/infrastructure/network"
	"holy-codex/infrastructure/storage"
	"holy-codex/services"

	//"github.com/libp2p/go-libp2p/config"
)

type Context struct {
	Config *config.Config
	Storage storage.Storage
	Discovery *network.Discovery
	Syncer *network.Sync
	AutoSave *services.AutoSave
	SyncService *services.SyncService
	Notifier *services.Notifier
}