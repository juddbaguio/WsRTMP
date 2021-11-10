package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/juddbaguio/WsRTMP"
)

var socketUpgrader = websocket.Upgrader{
	CheckOrigin: websocket.IsWebSocketUpgrade,
}

func main() {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/transcode/{streamKey}", func(w http.ResponseWriter, r *http.Request) {
		var mu sync.RWMutex
		c, err := socketUpgrader.Upgrade(w, r, nil)

		streamKey := mux.Vars(r)["streamKey"]

		rtmpDestination := fmt.Sprintf("rtmp://localhost:1935/live/%s", streamKey)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		rtmp := WsRTMP.New(rtmpDestination)

		err = rtmp.StartBroadcast()

		if err != nil {
			log.Panicf("failed to start RTMP broadcast: %w", err)
		}

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil || messageType == websocket.CloseMessage {
				mu.Lock()
				defer mu.Unlock()

				err := rtmp.StopBroadcast()
				if err != nil {
					log.Panic(fmt.Errorf("failed to kill FFmpeg process: %w", err))
				}
				c.Close()
				break
			}
			rtmp.StreamToPipe(&message)
		}
	})

	log.Fatal(http.ListenAndServe("localhost:3030", r))
}
