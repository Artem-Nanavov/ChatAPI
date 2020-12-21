package server

import (
	"api/entities"
	"api/utils"
	"encoding/json"

	"github.com/fasthttp/websocket"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// CreateChat ...
func (s *Server) CreateChat() routing.Handler {
	return func(c *routing.Context) error {
		var chat *entities.Chat
		if err := json.Unmarshal(c.Request.Body(), &chat); err != nil {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if err := s.db.Chat().Create(chat); err != nil {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return utils.Respond(c, 200, chat)
	}
}

// GetAllChatMessages ...
func (s *Server) GetAllChatMessages() routing.Handler {
	return func(c *routing.Context) error {
		id := c.Request.URI().QueryArgs().Peek("id")
		if id == nil {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": "id must be provided",
			})
		}
		messages, err := s.db.Message().GetAllByChatID(utils.ToInt(string(id)))
		if err != nil {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return utils.Respond(c, 200, messages)
	}
}

// Websocket ...
func (s *Server) Websocket() routing.Handler {
	return func(c *routing.Context) error {
		upgrader := websocket.FastHTTPUpgrader{}
		upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool { return true }
		upgrader.Upgrade(c.RequestCtx, func(ws *websocket.Conn) {
			// Request user goes online
			user := c.Get("user").(*entities.User)
			s.db.User().GoOnline(user)

			defer ws.Close()
			for {
				var message *entities.Message

				if err := ws.ReadJSON(&message); err != nil {
					ws.WriteJSON(map[string]interface{}{
						"error": err.Error(),
					})
					s.db.User().GoOfline(user)
					break
				}

				message.OwnerID = user.ID

				if err := s.db.Message().Create(message); err != nil {
					ws.WriteJSON(map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ws.WriteJSON(map[string]interface{}{
					"text":       message.Text,
					"created_at": message.CreatedAt,
					"username":   user.Username,
				})
			}
		})
		return nil
	}
}
