package utils

import (
	"encoding/json"

	routing "github.com/qiangxue/fasthttp-routing"
)

// Respond ...
func Respond(c *routing.Context, status int, data interface{}) error {
	c.SetStatusCode(status)
	return json.NewEncoder(c.Response.BodyWriter()).Encode(data)
}
