package server

import (
	"api/database"
	"api/services"

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

// GetRouting ...
func (s *Server) GetRouting() *routing.Router {
	router := routing.New()

	auth := router.Group("/auth")
	auth.Post("/reg", s.BaseMiddleware(s.RegistrationHandler()))
	auth.Post("/login", s.BaseMiddleware(s.LoginHandler()))

	user := router.Group("/users")
	user.Get("/me", s.BaseMiddleware(s.AuthenticationMiddleware(s.GetCurrentUser(), false)))
	user.Get("/", s.BaseMiddleware(s.AuthenticationMiddleware(s.GetAllUsers(), false)))

	chat := router.Group("/chats")
	chat.Post("/create", s.BaseMiddleware(s.AuthenticationMiddleware(s.CreateChat(), false)))
	chat.Get("/messages", s.BaseMiddleware(s.AuthenticationMiddleware(s.GetAllChatMessages(), false)))
	chat.Any("/ws", s.BaseMiddleware(s.AuthenticationMiddleware(s.Websocket(), true)))

	return router
}
