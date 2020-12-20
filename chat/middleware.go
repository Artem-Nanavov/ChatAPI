package chat

import (
	"api/entities"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func (w *WebSocketChat) parseUser(tokenString string) (*entities.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method, got " + t.Header["alg"].(string))
		}
		return []byte(w.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	return w.db.User().FindByID(int(id))
}
