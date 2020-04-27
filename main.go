package main

import (
	"flag"
	"log"

	"github.com/pkg/errors"

	"github.com/vllry/gameapi/pkg/backup"
	"github.com/vllry/gameapi/pkg/game"
	"github.com/vllry/gameapi/pkg/webserver"
)

func main() {
	var directory string
	var gcloudCredentialsPath string
	var instance string

	flag.StringVar(&directory, "directory", "", "Base directory for gameserver config and data.")
	flag.StringVar(&gcloudCredentialsPath, "gcloud-credentials-path", "", "Path to Google Cloud credentials file. No path will fall back to the magic SDK behavior.")
	flag.StringVar(&instance, "instance", "", "Game instance to run as. The instance represents a distinct world/save for a given game, e.g. different modpacks for Minecraft.")
	flag.Parse()

	if directory == "" {
		log.Fatal("Must supply --directory.")
	}

	if instance == "" {
		log.Fatal("Must supply --instance. Use --instance test for testing.")
	}

	backupManager := backup.NewManager(backup.NewGoogleCloudStorage(gcloudCredentialsPath), "/tmp")

	g, err := game.NewGame("minecraft", instance, directory, backupManager)
	if err != nil {
		panic(errors.Wrap(err, "this server has ants"))
	}

	// TODO handle graceful shutdown.
	webserver.Start(g)
}
