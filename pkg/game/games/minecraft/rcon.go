package minecraft

import (
	"errors"

	mcrcon "github.com/Kelwing/mc-rcon"
)

// Please give these things better names.

type rconCreatorIface interface {
	new() (rconClientIface, error)
}

type rconClientIface interface {
	run(string) (string, error)
	close()
}

type realRconConnection struct {
	conn *mcrcon.MCConn
}

type realRconCreator struct {
	port     int
	password string
}

func (r *realRconCreator) new() (rconClientIface, error) {
	conn := new(mcrcon.MCConn)

	err := conn.Open("", "")
	if err != nil {
		return nil, err
	}

	err = conn.Authenticate()
	if err != nil {
		return nil, err
	}

	return &realRconConnection{
		conn: conn,
	}, nil
}

func (r *realRconConnection) run(cmd string) (string, error) {
	return r.conn.SendCommand(cmd)
}

func (r *realRconConnection) close() {
	r.conn.Close()
}

type fakeRconCreator struct {
	// commands is a map of input commands to stdout results.
	commands map[string]string
}

func (r *fakeRconCreator) new() (rconClientIface, error) {
	return &fakeRconConnection{
		commands: r.commands,
	}, nil
}

type fakeRconConnection struct {
	commands map[string]string
}

func (r *fakeRconConnection) run(cmd string) (string, error) {
	output, found := r.commands[cmd]
	if found {
		return output, nil
	} else {
		return "", errors.New("unimplimented")
	}
}

func (r *fakeRconConnection) close() {}
