The Minecraft GameAPI sits as a daemon/sidecar adjacent to a Minecraft server
(EG same physical machine, or as a container in the same Kubernetes pod).

The GameAPI provides a standard, programmatic interface for the running game instance.

# Developing
Requires "recent" Go with GoModules support.

## Test
`go test ./...`

## Build
`go build`
