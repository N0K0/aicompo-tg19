package main

// Command struct defines how commands from players/admin should be formatted
type Command struct {
	Type  string `json:"type"`
	Value string `json:"command"`
}

type ClientInfo struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type StatusObject struct {
	NumPlayers int
	Players    map[string]Player
	GameStatus GameHandler
}
