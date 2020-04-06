package identifier

type GameIdentifier struct {
	// Game is the kind of game (e.g. Minecraft).
	Game string
	// Instance distinguishes individual servers of a game type (e.g. a vanilla Minecraft server vs a different modded one).
	Instance string
}
