package server

import (
	"api/entities"
	"api/utils"

	routing "github.com/qiangxue/fasthttp-routing"
)

// GetCurrentUser ...
func (s *Server) GetCurrentUser() routing.Handler {
	return func(c *routing.Context) error {
		user, ok := c.Get("user").(*entities.User)
		if !ok {
			return utils.Respond(c, 401, map[string]interface{}{
				"error": "Authentication credentials were not provided",
			})
		}
		return utils.Respond(c, 200, user)
	}
}
