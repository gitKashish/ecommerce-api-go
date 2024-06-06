package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gitKashish/EcomServer/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	subrouter := http.NewServeMux()
	subrouter.Handle("/v1/", http.StripPrefix("/v1", router))

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)
	fmt.Printf("Starting server at %s\n", s.addr)
	return http.ListenAndServe(s.addr, subrouter)
}
