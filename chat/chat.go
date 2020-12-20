package chat

import (
	"api/database"
	"api/entities"

	"github.com/fasthttp/websocket"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// Config ...
type Config struct {
	Port      string
	Logger    *logrus.Logger
	Salt      string
	JWTSecret string
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
			for {
				var message *entities.Message

				ws.SetCloseHandler(func(code int, text string) error {
					return ws.WriteJSON(map[string]interface{}{
						"connection": "closed",
					})
				})

				if err := ws.ReadJSON(&message); err != nil {
					ws.WriteJSON(map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				if err := w.db.Message().Create(message); err != nil {
					ws.WriteJSON(map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				user, _ := w.db.User().FindByID(message.OwnerID)

				ws.WriteJSON(map[string]interface{}{
					"text":       message.Text,
					"created_at": message.CreatedAt,
					"username":   user.Username,
				})
			}
		})
	}
}

// Run ...
func (w *WebSocketChat) Run() error {
	w.config.Logger.Info("Websocket server started at port " + w.config.Port)
	return fasthttp.ListenAndServe(":"+w.config.Port, w.handler())
}
