package websocketService

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=New

type WebsocketService struct {
	clients map[uuid.UUID]*Client
}

func New(service *WebsocketService) (*WebsocketService, error) {
	service.clients = make(map[uuid.UUID]*Client)
	return service, nil
}

func (service *WebsocketService) CreateConnection(
	conn *websocket.Conn,
	openHandler func(uuid.UUID) error,
	receiveHandler func(uuid.UUID, int, []byte) error,
	closeHandler func(uuid.UUID, error),
) error {
	if conn == nil || receiveHandler == nil {
		return errors.New("parameter is nil")
	}
	client := &Client{uuid.New(), conn}
	go func(client *Client) {
		defer conn.Close()
		defer delete(service.clients, client.id)
		for {
			err := client.Receive(receiveHandler)
			if err != nil {
				if closeHandler != nil {
					closeHandler(client.id, err)
				}
				break
			}
		}
	}(client)
	service.clients[client.id] = client
	return nil
}

func (service *WebsocketService) Send(id uuid.UUID, msgType int, msgBytes []byte) error {
	client, ok := service.clients[id]
	if !ok {
		return errors.New("client not existed")
	}
	if err := client.Send(msgType, msgBytes); err != nil {
		return err
	}
	return nil
}
