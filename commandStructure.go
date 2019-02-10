package main

// Command struct defines how commands from players/admin should be formatted
type Command struct {
	Type  string `json:"type"`
	Value string `json:"command"`
}

type Envelope struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ConfigUpdate struct {
	Pairs []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
}

type StatusObject struct {
	NumPlayers int
	Players    map[string]bool
	GameStatus GameHandler
}
