package cart

import (
	"fmt"

	"github.com/gitKashish/EcomServer/types"
)

// Return a list of only cart items & provides some validation(mentioned below).
// Just a utility function `types.CartItem` already has ProductID.
func getCartItemIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))

	for i, item := range items {
		// Check for invalid quantity:
		// because zero or negative quantity items should not exist...
		// as they may cause absurd checkout calculations, leading to...
		// negative quantities & 0 cost order items.
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

// Create order record in DB.
// Takes in Handler as receiver beacuse it need acess to a lot of other stores & types.
func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	// General Flow:
	/*
		1. All the cart items are in stock?
		TRUE:
			1. Calculate the total price.
			2. Reduce each items quantity in the DB (uses UpdateProduct() method).
			3. Create the order record in DB.
			4. Create order items record for each cart item.
			5. return orderID & total price.
		FALSE:
			1. return orderID: 0 & total_price: 0. Along with error.
	*/

	// Product Map created for quick lookup to entire...
	// `types.Product` struct using only Product ID of cart Items.
	productMap := make(map[int]types.Product)

	// Initializing Product Map
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

		h.productStore.UpdateProduct(product) // Updating records in DB.
	}
	// Create the order.
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // TODO : Get from http.Request.
		// TODO : Maintain a table for each user to store multiple addresses.
	})

	if err != nil {
		return 0, 0, err
	}

	// Create the OrderItems. For each cart Item.
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

// Util function to calculate total_price. Utilizes Product map for quick lookup.
func calculateTotalPrice(items []types.CartItem, products map[int]types.Product) (total float64) {
	for _, item := range items {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

// Check stock and sanity of cart Items.
func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {

	// Empty cart cannot be processed.
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	// If Product does not exists.
	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please update your cart", item.ProductID)
		}

		// Not enough stock in inventory.
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the inventory in the quantity requested", product.Name)
		}
	}

	return nil
}
