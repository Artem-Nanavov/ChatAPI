package server

import (
	"api/database"
	"api/services"
	"api/utils"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// Logger ...
type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
}

// Config ...
type Config struct {
	Port      string
	Logger    Logger
	Salt      string
	JWTSecret string
}

// Server ...
type Server struct {
	db          *database.Database
	config      *Config
	userService *services.UserService
}

// NewServer ...
func NewServer(db *database.Database, config *Config) *Server {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}
	return &Server{
		db:     db,
		config: config,
	}
}

// Run ...
func (s *Server) Run() error {
	s.config.Logger.Info("Server started at port " + s.config.Port + ".")
	if err := fasthttp.ListenAndServe(":"+s.config.Port, s.GetRouting().HandleRequest); err != nil {
		s.config.Logger.Error(err.Error())
		return err
	}
	return nil
}

// User ...
func (s *Server) User() *services.UserService {
	if s.userService == nil {
		s.userService = services.NewUserService(s.db, &services.Config{
			Salt:      s.config.Salt,
			JWTSecret: s.config.JWTSecret,
		})
	}
	return s.userService
}

// AuthenticationMiddleware ...
func (s *Server) AuthenticationMiddleware(handler routing.Handler) routing.Handler {
	return func(c *routing.Context) error {
		head := string(c.Request.Header.Peek("Authorization"))

		splited := strings.Split(head, " ")
		if len(splited) != 2 {
			return utils.Respond(c, 401, map[string]interface{}{
				"error": "Authentication credentials were not provided",
			})
		}

		token, err := jwt.Parse(splited[1], func(t *jwt.Token) (interface{}, error) {
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

// GetRouting ...
func (s *Server) GetRouting() *routing.Router {
	router := routing.New()

	auth := router.Group("/auth")
	auth.Post("/reg", s.BaseMiddleware(s.RegistrationHandler()))
	auth.Post("/login", s.BaseMiddleware(s.LoginHandler()))

	user := router.Group("/users")
	user.Get("/me", s.BaseMiddleware(s.AuthenticationMiddleware(s.GetCurrentUser())))
	user.Get("/", s.BaseMiddleware(s.AuthenticationMiddleware(s.GetAllUsers())))

	chat := router.Group("/chats")
	chat.Post("/create", s.BaseMiddleware(s.AuthenticationMiddleware(s.CreateChat())))
	chat.Get("/messages", s.BaseMiddleware(s.AuthenticationMiddleware(s.GetAllChatMessages())))

	return router
}
