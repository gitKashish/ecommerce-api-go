package cart

import (
	"fmt"

	"github.com/gitKashish/EcomServer/types"
)

func getCartItemIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}

		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}
	// Check if a product is in Stock.
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	// Calculate the total price.
	totalPrice := calculateTotalPrice(items, productMap)

	// Reduce the quantity of product.
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}
	// Create the order.
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // TODO : Get from user.
	})

	if err != nil {
		return 0, 0, err
	}

	// Create the OrderItems.
	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func calculateTotalPrice(items []types.CartItem, products map[int]types.Product) (total float64) {
	for _, item := range items {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please update your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the inventory in the quantity requested", product.Name)
		}
	}

	return nil
}
