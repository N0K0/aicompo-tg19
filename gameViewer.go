package main

import (
	"time"

	"github.com/google/logger"

	"github.com/gorilla/websocket"
)

type gameViewer struct {
	conn *websocket.Conn

	gm *GameHandler

	qSend chan []byte
	qRecv chan []byte

	statusTicker *time.Ticker
	pingTicker   *time.Ticker
}

func (gv *gameViewer) writePump() {
	logger.Info("Running gameViewer writePump")

	statusFreq := (100 / 10) * time.Millisecond

	logger.Infof("Update time: %v", statusFreq)
	gv.pingTicker = time.NewTicker(pingPeriod)
	defer func() {
		gv.pingTicker.Stop()
		gv.statusTicker.Stop()
	}()
	for {
		select {
		case message, ok := <-gv.qSend:

			if gv.conn == nil {
				logger.Infof("No more connection, can not print")
				return
			}
			err := gv.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Info("SetWrDeadline err")
			}

			if !ok {
				// The hub closed the channel.
				err := gv.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Info("WriteMessage err")
				}
				return
			}

			w, err := gv.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Info("Error with NextWriter")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				logger.Info("Error write")
			}

			// Add queued chat messages to the current websocket message.
			n := len(gv.qSend)
			for i := 0; i < n; i++ {
				_, err := w.Write(<-gv.qSend)
				if err != nil {
					logger.Info("Error write")
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-gv.pingTicker.C:
			//logger.Info("ping ticker")
			err := gv.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Info("Error set deadline Ticker")
			}
			if err := gv.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}

	}
}

func (gv *gameViewer) readPump() {
	logger.Info("Running gameViewer readPump")
	err := gv.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logger.Info("Err set read Deadline")
	}
	gv.conn.SetPongHandler(
		func(string) error {
			err := gv.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				logger.Infof("Err set deadline pong")
			}
			return nil
		},
	)

	for {
		_, message, err := gv.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				logger.Infof("IsUnexpectedCloseError GameViewer: %v", err)
			} else {
				logger.Infof("GameViewer closed socket at %v ", gv.conn.RemoteAddr())
				gv.conn = nil
				gv.pingTicker.Stop()
			}
			break
		}
		gv.qRecv <- message
	}
}

func (gv *gameViewer) statusUpdater() {
	logger.Info("Running statusUpdater")
	gv.statusTicker = time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-gv.statusTicker.C:
			logger.Info("Pushing new status")
			gv.qSend <- gv.gm.generateStatusJson()
		}
	}
}

func (gv *gameViewer) run() {
	logger.Info("Running viewer")
	go gv.readPump()
	go gv.writePump()
	go gv.statusUpdater()
}
