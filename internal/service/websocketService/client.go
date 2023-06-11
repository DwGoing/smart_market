package websocketService

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id   uuid.UUID
	conn *websocket.Conn
}

func (client *Client) Receive(handler func(uuid.UUID, int, []byte) error) error {
	msgType, msgBytes, err := client.conn.ReadMessage()
	if err != nil {
		return err
	}
	err = handler(client.id, msgType, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Send(msgType int, msgBytes []byte) error {
	err := client.conn.WriteMessage(msgType, msgBytes)
	return err
}
