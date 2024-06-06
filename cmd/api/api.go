package api

import (
	"database/sql"
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

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(router)

	return http.ListenAndServe(s.addr, subrouter)
}
