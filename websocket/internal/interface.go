package websocket

import "github.com/gorilla/websocket"

type IHub interface {
	AddClient(client *websocket.Conn)
	Run()
}
