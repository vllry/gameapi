The Minecraft GameAPI sits as a daemon/sidecar adjacent to a Minecraft server
(EG same physical machine, or as a container in the same Kubernetes pod).

The GameAPI provides a standard, programmatic interface for the running game instance.

# Developing
* Requires "recent" Go with GoModules support.
* `go get golang.org/x/tools/cmd/goimports`

## Test
`make test`

## Build
`make build` (or `go build`)

## General structure
`./pkg/webserver` contains a basic HTTP webserver,
and API wrappers around Game object calls.

`./pkg/game` contains the generic constructor for new Games.

`.pkg/game/gameinterface` contains the interface definition for Games.

`./pkg/game/games` contains individual Game implementations,
which satisfy the generic Game interface.
