package main

// Command struct defines how commands from players/admin should be formatted
type Command struct {
	CommandType string `json:"type"`
	Command     string `json:"command"`
}
