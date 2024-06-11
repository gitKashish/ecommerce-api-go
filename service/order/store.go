package order

import (
	"database/sql"

	"github.com/gitKashish/ecommerce-api-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// create a new order record in `orders` table in DB...
// and return the order ID.
func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	// Getting the last queries `AUTO_INCREMENT` generated ID.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Create a new order item record in `order_items` table in DB.
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}
