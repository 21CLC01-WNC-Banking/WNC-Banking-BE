package controller

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/websocket"
)

type ApiContainer struct {
	HttpServer      *http.Server
	WebsocketServer *websocket.Server
}

func NewApiContainer(httpServer *http.Server, websocketServer *websocket.Server) *ApiContainer {
	return &ApiContainer{HttpServer: httpServer, WebsocketServer: websocketServer}
}
