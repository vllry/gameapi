The GameAPI sits as a local facade for a supported game.

The GameAPI provides a standard, programmatic interface for the running game instance.
The objective of the GameAPI is to abstract RCON and filesystem operations into a safer (and easy to consume) API.

# Dependencies
* 7zr
* Google Cloud Storage credentials + access.

# Developing
* Requires "recent" Go with GoModules support.
* goimpors tool: `go get golang.org/x/tools/cmd/goimports`

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
