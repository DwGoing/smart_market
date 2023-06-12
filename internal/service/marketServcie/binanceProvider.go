package marketServcie

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewBinanceProvider

type BinanceProvider struct {
	connection *websocket.Conn
	requests   map[string]Request
}

func NewBinanceProvider(provider *BinanceProvider) (*BinanceProvider, error) {
	log.SetPrefix("[BinanceProvider] ")
	provider.initialize()
	provider.requests = make(map[string]Request)
	return provider, nil
}

func (provider *BinanceProvider) initialize() {
	conn, _, err := websocket.DefaultDialer.Dial("wss://ws-api.binance.com:443/ws-api/v3", nil)
	if err != nil {
		log.Printf("initializ error: %s", err)
		time.Sleep(time.Second * 3)
		provider.initialize()
		return
	}
	provider.connection = conn
	go func() {
		defer conn.Close()
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Printf("eceive message error: %s", err)
				provider.initialize()
				return
			} else {
				log.Printf("receive message: %s", msgBytes)
				err = provider.handleMessage(msgBytes)
				if err != nil {
					log.Printf("handle message error: %s | %s", msgBytes, err)
				}
			}
		}
	}()
	// health check
	go func() {
		request := Request{Id: uuid.NewString(), Method: "ping"}
		requestBytes, _ := json.Marshal(request)
		for {
			err := conn.WriteMessage(1, requestBytes)
			if err != nil {
				log.Printf("ping error: %s", err)
				provider.initialize()
				return
			}
			time.Sleep(time.Second * 30)
		}
	}()
}

func (provider *BinanceProvider) handleMessage(msgBytes []byte) error {
	var response Request
	err := json.Unmarshal(msgBytes, &response)
	if err != nil {
		return err
	}
	request, ok := provider.requests[response.Id]
	if ok {
		switch request.Method {
		case "ping":
			log.Printf("XXX")
		}
		delete(provider.requests, response.Id)
	}
	return nil
}
