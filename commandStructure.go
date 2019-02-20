package main

import "encoding/json"

// Command struct defines how commands from players/admin should be formatted
type Command struct {
	Type  string `json:"type"`
	Value string `json:"command"`
}

type Envelope struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type EnvelopePartial struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

type ConfigUpdate struct {
	Configs []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"configs"`
}

type StatusObject struct {
	NumPlayers int
	Players    map[string]bool
	GameStatus GameHandler
}
