package main

// Command struct defines how commands from players/admin should be formatted
type Command struct {
	Command string `json:"command"`
}

// Structure defining game status to the players
type gameStatus struct {
	updateType string
	data       string
}
