package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gitKashish/EcomServer/service/cart"
	"github.com/gitKashish/EcomServer/service/order"
	"github.com/gitKashish/EcomServer/service/product"
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

	// User handler service
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	// Product handler service
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(router)

	// Cart handler service
	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, userStore, productStore)
	cartHandler.RegisterRoutes(router)

	fmt.Printf("Starting server at %s\n", s.addr)
	return http.ListenAndServe(s.addr, subrouter)
}
