package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const ()

type gameViewer struct {
	conn *websocket.Conn

	gm *GameHandler

	qSend chan []byte
	qRecv chan []byte

	statusTicker *time.Ticker
	pingTicker   *time.Ticker
}

func (gv *gameViewer) writeSocket() {

	statusFreq := (1000 / 10) * time.Millisecond

	log.Printf("Update time: %v", statusFreq)
	gv.pingTicker = time.NewTicker(pingPeriod)
	gv.statusTicker = time.NewTicker(statusFreq)
	defer func() {
		gv.pingTicker.Stop()
		gv.statusTicker.Stop()
	}()
	for {
		select {
		case message, ok := <-gv.qSend:

			if gv.conn == nil {
				log.Printf("No more connection, can not print")
				return
			}
			err := gv.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("SetWrDeadline err")
			}

			if !ok {
				// The hub closed the channel.
				err := gv.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Print("WriteMessage err")
				}
				return
			}

			w, err := gv.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print("Error with NextWriter")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Print("Error write")
			}

			// Add queued chat messages to the current websocket message.
			n := len(gv.qSend)
			for i := 0; i < n; i++ {
				_, err := w.Write(<-gv.qSend)
				if err != nil {
					log.Print("Error write")
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-gv.pingTicker.C:
			err := gv.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("Error set deadline Ticker")
			}
			if err := gv.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-gv.statusTicker.C:
			// Time to update the status
			gv.gm.mapLock.Lock()
			if gv.conn == nil {
				gv.statusTicker.Stop()
			}
			gv.qSend <- gv.gm.generateStatusJson()
			gv.gm.mapLock.Unlock()

		}
	}
}

func (gv *gameViewer) readSocket() {
	log.Print("GameViewer Read Socket")
	gv.conn.SetReadLimit(maxMessageSize)
	err := gv.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Print("Err set read Deadline")
	}
	gv.conn.SetPongHandler(
		func(string) error {
			err := gv.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				log.Printf("Err set deadline pong")
			}
			return nil
		},
	)

	for {
		_, message, err := gv.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("IsUnexpectedCloseError GameViewer: %v", err)
			} else {
				log.Printf("GameViewer closed socket at %v ", gv.conn.RemoteAddr())
				gv.conn = nil
				gv.pingTicker.Stop()
			}
			break
		}
		log.Print("Waiting for GameViewer message")
		gv.qRecv <- message
	}
}

func (gv *gameViewer) run() {
	go gv.readSocket()
	go gv.writeSocket()

}
