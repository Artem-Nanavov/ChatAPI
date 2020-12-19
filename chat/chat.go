package chat

import (
	"api/database"
	"api/entities"
	"encoding/json"

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
		upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				messageType, p, _ := ws.ReadMessage()

				var message *entities.Message
				if err := json.Unmarshal(p, &message); err != nil {
					data, _ := json.Marshal(map[string]interface{}{
						"error": "Cannot serialize data",
					})
					ws.WriteMessage(messageType, data)
					continue
				}

				if err := w.db.Message().Create(message); err != nil {
					data, _ := json.Marshal(map[string]interface{}{
						"error": err.Error(),
					})
					ws.WriteMessage(messageType, data)
					continue
				}

				data, _ := json.Marshal(message)
				ws.WriteMessage(messageType, data)
			}
		})
	}
}

// Run ...
func (w *WebSocketChat) Run() error {
	w.config.Logger.Info("Websocket server started at port " + w.config.Port)
	return fasthttp.ListenAndServe(":"+w.config.Port, w.handler())
}
