package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/DwGoing/smart_market/internal/service/websocketService"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	WebsocketService *websocketService.WebsocketService `singleton:""`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *App) Run() {
	engine := gin.Default()
	engine.GET("/ws", app.createWebsocketConnection)
	_ = engine.Run("0.0.0.0:9000")
}

func (app *App) createWebsocketConnection(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Request.Trailer)
	if err != nil {
		ctx.Writer.Write([]byte(err.Error()))
	}
	err = app.WebsocketService.CreateConnection(
		conn,
		nil,
		func(id uuid.UUID, msgType int, msgBytes []byte) error {
			if err = app.handleWebsocketMessage(id, msgType, msgBytes); err != nil {
				log.Printf("handle websocket message error: %s | %s | %s", id, err, msgBytes)
			}
			return nil
		},
		nil,
	)
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
	}
}

func (app *App) handleWebsocketMessage(id uuid.UUID, msgType int, msgBytes []byte) error {
	log.Printf("receive websocket message: %s | %s", id, msgBytes)
	var request Request
	err := json.Unmarshal(msgBytes, &request)
	if err != nil {
		return err
	}
	switch request.Method {
	case "ping":
		responseBytes, err := json.Marshal(Response{Id: request.Id, Code: 200, Message: "pong"})
		if err != nil {
			return err
		}
		err = app.WebsocketService.Send(id, 1, responseBytes)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown request")
	}
	return nil
}
