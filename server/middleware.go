package server

import (
	"api/utils"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
)

// AuthenticationMiddleware ...
func (s *Server) AuthenticationMiddleware(handler routing.Handler, ws bool) routing.Handler {
	return func(c *routing.Context) error {
		var unprocessedToken string
		if ws {
			unprocessedToken = string(c.QueryArgs().Peek("token"))
		} else {
			head := string(c.Request.Header.Peek("Authorization"))
			splited := strings.Split(head, " ")
			if len(splited) != 2 {
				return utils.Respond(c, 401, map[string]interface{}{
					"error": "Authentication credentials were not provided",
				})
			}
			unprocessedToken = splited[1]
		}

		token, err := jwt.Parse(unprocessedToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method, got " + t.Header["alg"].(string))
			}
			return []byte(s.config.JWTSecret), nil
		})

		if err != nil {
			return utils.Respond(c, 401, map[string]interface{}{
				"error": err.Error(),
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"].(float64)
		user, err := s.db.User().FindByID(int(id))
		if err != nil {
			return utils.Respond(c, 401, map[string]interface{}{
				"error": err.Error(),
			})
		}

		c.Set("user", user)
		handler(c)
		return nil
	}
}

// BaseMiddleware ...
func (s *Server) BaseMiddleware(handler routing.Handler) routing.Handler {
	return func(c *routing.Context) error {
		c.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		c.Response.Header.SetBytesV("Access-Control-Allow-Origin", c.Request.Header.Peek("Origin"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler(c)

		s.config.Logger.Info(
			string(string(c.Request.URI().Path()) + " - " + string(c.Request.Header.Method())),
		)
		return nil
	}
}
