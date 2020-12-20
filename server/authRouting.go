package server

import (
	"api/entities"
	"api/utils"
	"database/sql"
	"encoding/json"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// RegistrationHandler ...
func (s *Server) RegistrationHandler() routing.Handler {
	return func(c *routing.Context) error {

		var user entities.User
		json.Unmarshal(c.Request.Body(), &user)

		token, err := s.User().Create(&user)
		if err != nil {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": err.Error(),
			})
		}

		cookie := fasthttp.Cookie{}
		cookie.SetKey("token")
		cookie.SetValue(token)
		cookie.SetHTTPOnly(true)
		cookie.SetPath("/")
		c.Response.Header.SetCookie(&cookie)

		return utils.Respond(c, 200, map[string]interface{}{
			"token":    token,
			"username": user.Username,
			"id":       user.ID,
		})
	}
}

// LoginHandler ...
func (s *Server) LoginHandler() routing.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(c *routing.Context) error {
		var data request
		json.Unmarshal(c.Request.Body(), &data)

		user, err := s.User().Repo().FindByEmail(data.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return utils.Respond(c, 400, map[string]interface{}{
					"error": "No user with such email: " + data.Email,
				})
			}
			return utils.Respond(c, 400, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if !s.User().ComparePasswords(user, data.Password) {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": "Wrong email or password",
			})
		}

		token, err := s.User().GenerateToken(user)
		if err != nil {
			return utils.Respond(c, 400, map[string]interface{}{
				"error": err.Error(),
			})
		}

		cookie := fasthttp.Cookie{}
		cookie.SetKey("token")
		cookie.SetValue(token)
		cookie.SetHTTPOnly(true)
		cookie.SetPath("/")
		c.Response.Header.SetCookie(&cookie)

		return utils.Respond(c, 200, map[string]interface{}{
			"token":    token,
			"username": user.Username,
			"id":       user.ID,
		})
	}
}
