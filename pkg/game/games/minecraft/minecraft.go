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
		rconConstructor: &realRcon{ // Default to using real RCON.
			password: rconPassword,
			port: rconPort,
		},
	}

	return &Game{
		config: config,
	}, nil
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
	for _, player := range strings.Split(playersString, ",") {
		players = append(players, player)
	}

	return players, nil
}
