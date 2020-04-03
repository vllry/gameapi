package minecraft

import mcrcon "github.com/Kelwing/mc-rcon"

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

type realRcon struct {
	port int
	password string
}

func (r *realRcon) new() (rconClientIface, error) {
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
