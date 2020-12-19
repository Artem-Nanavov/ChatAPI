package chat

import (
	"api/database"
	"api/entities"
	"encoding/json"
	"fmt"

	"github.com/fasthttp/websocket"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// Config ...
type Config struct {
	Port   string
	Logger *logrus.Logger
}

// WebSocketChat ...
type WebSocketChat struct {
	db     *database.Database
	config *Config
}

// NewWebSocketChat ...
func NewWebSocketChat(db *database.Database, config *Config) *WebSocketChat {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}
	return &WebSocketChat{
		db:     db,
		config: config,
	}
}

func (w *WebSocketChat) handler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		upgrader := websocket.FastHTTPUpgrader{}

		upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool { return true }
		upgrader.Upgrade(ctx, connection)
	}
}

func connection(ws *websocket.Conn) {
	defer ws.Close()
	for {
		messageType, p, _ := ws.ReadMessage()

		var message *entities.Message
		json.Unmarshal(p, &message)

		newMessage := string(message.Text)
		fmt.Println(newMessage)
		
		ws.WriteMessage(messageType, []byte(newMessage))
	}
}

// Run ...
func (w *WebSocketChat) Run() error {
	w.config.Logger.Info("Websocket server started at port " + w.config.Port)
	return fasthttp.ListenAndServe(":"+w.config.Port, w.handler())
}
