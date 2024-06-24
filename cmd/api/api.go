package api

import (
	"database/sql"
	"log"

	"github.com/dclouisDan/chat-app-api/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST,PUT,DELETE,GET,OPTIONS",
		AllowHeaders: "Origin, Accept, Content-Type, Authorization",
	}))

  api := app.Group("/chat-app-api/v1")
  api.Static("/", "./web/static")

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(api)
  
	log.Println("Listening on:", s.addr)
  return app.Listen(s.addr)
}
