package main

import (
	"flag"

	"github.com/pkg/errors"

	"github.com/vllry/gameapi/pkg/backup"
	"github.com/vllry/gameapi/pkg/game"
	"github.com/vllry/gameapi/pkg/webserver"
)

func main() {
	var gcloudCredentialsPath string
	flag.StringVar(&gcloudCredentialsPath, "gcloud-credentials-path", "", "Path to Google Cloud credentials file. No path will fall back to the magic SDK behavior.")
	flag.Parse()

	backupManager := backup.NewManager(backup.NewGoogleCloudStorage(gcloudCredentialsPath))

	// TODO take input.
	g, err := game.NewGame("minecraft", "fake", "test/minecraft", backupManager)
	if err != nil {
		panic(errors.Wrap(err, "this server has ants"))
	}

	// TODO handle graceful shutdown.
	webserver.Start(g)
}
