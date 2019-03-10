package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

// adminHandler takes care of starting the game and is used for spectating the events. Used by the web interface
type adminHandler struct {
	gm   *GameHandler
	conn *websocket.Conn

	qSend chan []byte
	qRecv chan []byte

	ticker *time.Ticker
}

func (admin *adminHandler) run() {
	logger.Infof("Admin loop started")
	logger.Infof("Game state: %v", admin.gm.Status)

	go admin.readSocket()
	go admin.writeSocket()

	loggerTicker := time.NewTicker(10 * time.Second)
	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case incoming := <-admin.qRecv:
			logger.Infof("New admin message: %s", incoming)
			admin.adminParseCommand(incoming)
			break
		case <-loggerTicker.C:
			admin.logStatus()
			break
		}
	}
}

func (admin *adminHandler) adminParseCommand(jsonObj []byte) {
	logger.Info("Got admin command")
	c := EnvelopePartial{}

	err := json.Unmarshal(jsonObj, &c)

	if err != nil {
		logger.Infof("Problem with admin command %v: %v", c, err)
	}

	switch c.Type {
	case "config":
		admin.adminParseConfigUpdates(&c.Message)
	case "config_get":
		admin.adminPushConfig()
	case "game_start":
		break
	default:
		break
	}

}

func (admin *adminHandler) adminParseConfigUpdates(jsonObj *json.RawMessage) {
	logger.Infof("Admin Config update\n%v", jsonObj)
	var c ConfigUpdate
	err := json.Unmarshal(*jsonObj, &c)

	if err != nil {
		logger.Infof("Problem with admin command %v: %v", c, err)
	}

	logger.Infof("Got json config: %v", c)

	ch := admin.gm.config

	for _, config := range c.Configs {
		switch config.Name {
		case "minTurnUpdate":
			if ch.MinTurnUpdate > ch.MaxTurnUpdate {
				admin.sendError("Min turn time is not lower than max turn time!")
				break
			}

			value, err := strconv.Atoi(config.Value)
			if err != nil {
				admin.sendError("Unable to convert value to int")
			}
			ch.MinTurnUpdate = value

		case "maxTurnUpdate":
			if ch.MinTurnUpdate > ch.MaxTurnUpdate {
				admin.sendError("Max turn time is lower than min turn time!")
				break
			}

			value, err := strconv.Atoi(config.Value)
			if err != nil {
				admin.sendError("Unable to convert value to int")
			}
			ch.MaxTurnUpdate = value

		case "mapSize":
			parts := strings.Split(config.Value, "x")

			if len(parts) != 2 {
				admin.sendError("Did not find two elements when splitting on 'X'")
				break
			}

			mapXStr, mapYStr := parts[0], parts[1]

			mapX, err := strconv.Atoi(mapXStr)
			if err != nil {
				admin.sendError("Could not convert X elem to int")
				break
			}

			if mapX < 3 && mapX != 0 {
				admin.sendError("Map can not be under size 3")
				break
			}
			mapY, err := strconv.Atoi(mapYStr)
			if err != nil {
				admin.sendError("Could not convert Y elem to int")
				break
			}

			if mapY < 3 && mapY != 0 {
				admin.sendError("Map can not be under size 3")
				break
			}
			ch.mapSizeX = mapX
			ch.mapSizeY = mapY

		case "outerWalls":

		default:
			logger.Error("Unable to parse configLine %v", config.Name)
		}
	}

}

func (admin *adminHandler) adminPushConfig() {
	logger.Info("Pushing config to frontend")

	jsonString, err := json.Marshal(admin.gm.config)
	if err != nil {
		admin.sendError("Problems with creating message")
		return
	}

	env := Envelope{
		Type:    "config_push",
		Message: string(jsonString),
	}

	jsonString, err = json.Marshal(env)
	if err != nil {
		admin.sendError("Problems with creating message")
		return
	}

	logger.Info(string(jsonString))

	admin.qSend <- jsonString
}

func (admin *adminHandler) logStatus() {
	statusString := `Admin status:
		Connected Players: %v
	`
	numPlayers := len(admin.gm.players)
	logger.Infof(statusString, numPlayers)
}

func (admin *adminHandler) sendError(message string) {
	msg := Envelope{
		Type:    "error",
		Message: message,
	}

	jsonString, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Problems with creating error message")
		return
	}

	logger.Errorf("Sending error: %s", message)
	admin.qSend <- jsonString
}

func (admin *adminHandler) writeSocket() {
	logger.Info("Write Socket")
	admin.ticker = time.NewTicker(pingPeriod)
	defer func() {
		admin.ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-admin.qSend:
			err := admin.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Info("SetWrDeadline err")
			}

			if !ok {
				// The hub closed the channel.
				err := admin.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Info("WriteMessage err")
				}
				return
			}

			w, err := admin.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Info("Error with NextWriter")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				logger.Info("Error write")
			}

			// Add queued chat messages to the current websocket message.
			n := len(admin.qSend)
			for i := 0; i < n; i++ {
				_, err := w.Write(<-admin.qSend)
				if err != nil {
					logger.Info("Error write")
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-admin.ticker.C:
			logger.Info("Admin ping ticker")
			err := admin.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Info("Error set deadline Ticker")
			}
			if err := admin.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (admin *adminHandler) readSocket() {
	logger.Info("Admin Read Socket")
	err := admin.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logger.Info("Err set read Deadline")
	}
	admin.conn.SetPongHandler(
		func(string) error {
			err := admin.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				logger.Infof("Err set deadline pong")
			}
			return nil
		},
	)

	for {
		_, message, err := admin.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				logger.Infof("IsUnexpectedCloseError admin: %v", err)
			} else {
				// Admin disconnected
				logger.Infof("Admin closed socket at %v . Removing connection", admin.conn.RemoteAddr())
				admin.conn = nil
				admin.ticker.Stop()
			}
			break
		}
		logger.Info("Waiting for admin message")
		admin.qRecv <- message
	}
}
