package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(hub IHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("HandleWebSocket: Failed to upgrade connection:", err)
			return
		}

		hub.AddClient(conn)
		fmt.Println("HandleWebSocket: Client added to Hub")
	}
}
