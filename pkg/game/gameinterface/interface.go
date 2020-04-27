package gameinterface

// This package must be widely imported.
// Therefore, be cautious of import loops when importing module to this module.

import (
	"github.com/vllry/gameapi/pkg/backup"
	"github.com/vllry/gameapi/pkg/game/identifier"
)

// Config contains general information needed by all game instances (EG what the instance name is).
type Config struct {
	BackupManager *backup.Manager
	Identifier    identifier.GameIdentifier
	GameDirectory string
}

// GenericGame provides standard functions for an arbitrary game.
type GenericGame interface {
	Backup() error
	GetLogs() (string, error)
	ListPlayers() ([]string, error)
}
