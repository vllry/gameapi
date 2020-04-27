package minecraft

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

// skipBackup is a set of file names to not back up.
var skipBackup = map[string]struct{}{
	"backups": {}, // FTB Revelation creates local backups by default. It's overkill to capture those duplicates.
}

// Game is the state of the Minecraft server manager.
type Game struct {
	config Config
}

// Config is Minecraft-specific configuration.
type Config struct {
	base            gameinterface.Config
	rconPort        int
	rconPassword    string
	rconConstructor rconCreatorIface // Mockable interface for starting realRconConnection sessions.
}

// NewGame initializes and returns a new Minecraft Game instance.
func NewGame(baseConfig gameinterface.Config) (gameinterface.GenericGame, error) {
	// Load the server.properties file.
	f, err := os.Open(path.Join(baseConfig.GameDirectory, "server.properties"))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	serverPropertiesBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// Fetch RCON settings.
	rconPort := 0
	rconPassword := ""
	for _, line := range strings.Split(string(serverPropertiesBytes), "\n") {
		if strings.HasPrefix(line, "rcon.port=") {
			rconPort, err = strconv.Atoi(strings.TrimLeft(line, "rcon.port="))
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "rcon.password=") {
			rconPassword = strings.TrimPrefix(line, "rcon.password=")
		}
	}
	if rconPort == 0 || len(rconPassword) == 0 {
		return nil, errors.New("unable to load rcon settings")
	}

	// Build the config object.
	config := Config{
		base: baseConfig,
		rconConstructor: &realRconCreator{ // Default to using real RCON.
			password: rconPassword,
			port:     rconPort,
		},
	}

	game := buildGame(config)
	return game, nil
}

func buildGame(config Config) *Game {
	return &Game{
		config: config,
	}
}

// Backup backs up the game.
// Only 1 copy of a game instance (distinct game + name pair) should be active at a time.
func (g *Game) Backup() error {
	rcon, err := g.config.rconConstructor.new()
	if err != nil {
		return errors.Wrap(err, "couldn't create rcon client")
	}
	defer rcon.close()
	defer rcon.run("save-on") // Always re-enable autosaving on exit.

	// Disable autosaving.
	_, err = rcon.run("save-off")
	if err != nil {
		return errors.Wrap(err, "couldn't disable auosaving")
	}

	// Manually flush save to disk.
	_, err = rcon.run("save")
	if err != nil {
		return errors.Wrap(err, "couldn't save game")
	}

	// Archive directory while saving is off.
	log.Println("Starting archive...")
	backupFiles := make([]string, 0)
	allFiles, err := ioutil.ReadDir(g.config.base.GameDirectory)
	for _, file := range allFiles {
		if _, found := skipBackup[file.Name()]; !found {
			backupFiles = append(backupFiles, file.Name())
		}
	}

	backup, err := g.config.base.BackupManager.ArchiveDirectory(g.config.base.GameDirectory, backupFiles, g.config.base.Identifier)
	if err != nil {
		return errors.Wrap(err, "failed to archive game directory")
	}

	// Re-enable saving, we're done disk IO.
	log.Println("Archiving done, re-enabling saving...")
	_, err = rcon.run("save-on")
	if err != nil {
		return errors.Wrap(err, "failed to re-enable saving") // TODO don't bail out here. May as well upload.
	}

	// upload archive.
	log.Println("Uploading...")
	err = backup.Upload(context.Background())
	return err
}

func (g *Game) GetLogs() (string, error) {
	// TODO We're trusting here that the log is rotated reasonably...
	// TODO Should older log files be spliced together?
	bytes, err := ioutil.ReadFile(path.Join(g.config.base.GameDirectory, "logs/latest.log"))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (g *Game) ListPlayers() ([]string, error) {
	rc, err := g.config.rconConstructor.new()
	if err != nil {
		return nil, err
	}
	defer rc.close()

	resp, err := rc.run("list")
	if err != nil {
		log.Fatalln("Command failed", err)
	}

	// TODO handle format errors.
	players := make([]string, 0)
	playersString := strings.Split(resp, "players online:")[1]
	if len(playersString) == 0 {
		return players, nil
	}

	for _, player := range strings.Split(playersString, ",") {
		players = append(players, player)
	}

	return players, nil
}
