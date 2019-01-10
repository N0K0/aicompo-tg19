package main

// Command struct defines how commands from players/admin should be formatted
type Command struct {
	command string
}

// Structure defining gamestatus to the players
type gameStatus struct {
	updateType string
	data       string
}
