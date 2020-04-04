package minecraft

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/vllry/gameapi/pkg/game/gameinterface"
)

type Game struct {
	config Config
}

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
			rconPassword = strings.TrimLeft(line, "rcon.password=")
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

// TODO (longterm) Backups should eventually be streamed to a central agent, rather than placed in an object store directly.
// Only 1 copy of a game instance (distinct game + name pair) should be active at a time.
// A rogue GameAPI instance should never be able to pollute the list of backups with a non-fresh backup.
func (g *Game) Backup() error {
	return nil
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
