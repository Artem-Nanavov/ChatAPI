package server

import (
	"api/entities"
	"api/utils"
	"encoding/json"

	routing "github.com/qiangxue/fasthttp-routing"
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

		return utils.Respond(c, 400, chat)
	}
}