package main

import (
	"github.com/pkg/errors"

	"github.com/vllry/gameapi/pkg/webserver"

	"github.com/vllry/gameapi/pkg/game"
)

func main() {
	// TODO take input.
	g, err := game.NewGame("minecraft", "default", "test/minecraft")
	if err != nil {
		panic(errors.Wrap(err, "this server has ants"))
	}

	// TODO handle graceful shutdown.
	webserver.Start(g)
}
