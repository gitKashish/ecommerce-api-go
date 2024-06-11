package cart

import (
	"fmt"
	"net/http"

	"github.com/gitKashish/ecommerce-api-go/service/auth"
	"github.com/gitKashish/ecommerce-api-go/types"
	"github.com/gitKashish/ecommerce-api-go/utils"
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

// Handler Functions for performing checkout operations.
func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUseIDFromContext(r.Context()) // Getting userID from current context.
	// This method is invoked by `auth.WithJWTAuth()` which updates Context...
	// before invoking.

	var cart types.CartCheckoutPayload
	// Parsing JSON request body to `types.CartCheckoutPayload`.
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validating Payload structure.
	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload : %v", errors))
		return
	}

	// Getting only productIDs from cart items.
	productIDs, err := getCartItemIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Getting list of `types.Product` for product details of cart items from DB.
	ps, err := h.productStore.GetProductByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Creating order record in `orders` table.
	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Responding on successful checkout.
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}
