package cart

import (
	"fmt"
	"net/http"

	"github.com/gitKashish/EcomServer/service/auth"
	"github.com/gitKashish/EcomServer/types"
	"github.com/gitKashish/EcomServer/utils"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	orderStore   types.OrderStore
	userStore    types.UserStore
	productStore types.ProductStore
}

func NewHandler(orderStore types.OrderStore, userStore types.UserStore, productStore types.ProductStore) *Handler {
	return &Handler{orderStore: orderStore, userStore: userStore, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore))
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUseIDFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload : %v", errors))
		return
	}

	// get products (slice)
	productIDs, err := getCartItemIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}
