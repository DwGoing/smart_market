package app

import (
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
	engine.GET("/ws", app.ws)
	_ = engine.Run("0.0.0.0:9000")
}

func (app *App) ws(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Request.Trailer)
	if err != nil {
		ctx.Writer.Write([]byte(err.Error()))
	}
	err = app.WebsocketService.CreateConnection(conn, func(id uuid.UUID, msgType int, msgBytes []byte) error {
		log.Printf("reveive msg: %s %s", msgBytes, id)
		app.WebsocketService.Send(id, msgType, msgBytes)
		return nil
	}, func(id uuid.UUID, err error) {
		log.Printf("%s exited %s", id, err)
	})
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
	}
}
