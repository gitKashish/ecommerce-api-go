package product

import (
	"net/http"

	"github.com/gitKashish/EcomServer/types"
	"github.com/gitKashish/EcomServer/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /products", h.handleGetProducts)
}

// HandlerFunc to get products (list)
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// General Flow:
	// 1. Get list of products from DB.
	// 2. Write a JSON response.

	// Getting a list of products from DB.
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}
