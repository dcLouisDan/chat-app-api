package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dclouisDan/chat-app-api/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct{
  addr string
  db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
  return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
  router := mux.NewRouter()
  subrouter := router.PathPrefix("/chat-app-api/v1").Subrouter()
  
  userStore := user.NewStore(s.db)
  userHandler := user.NewHandler(userStore)
  userHandler.RegisterRoutes(subrouter)

  log.Println("Listening on:", s.addr)
  return http.ListenAndServe(s.addr, router)
}
