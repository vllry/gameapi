package gameinterface

// Config contains general information needed by all game instances (EG what the instance name is).
type Config struct {
	InstanceName  string
	GameDirectory string
}

// GenericGame provides standard functions for an arbitrary game.
type GenericGame interface {
	GetLogs() (string, error)
	ListPlayers() ([]string, error)
}
