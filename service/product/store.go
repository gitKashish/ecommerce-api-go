package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gitKashish/ecommerce-api-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Get a list of products currently in the inventory.
func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products") // SELECT Query.
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	for rows.Next() {
		// transforming fetched record(row) into `types.Product`.
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

// Return a list `types.Product` from a list of Product IDs.
func (s *Store) GetProductByIDs(productIDs []int) ([]types.Product, error) {

	// Creating a query to select product records...
	// of products with given product ID.
	placeholders := strings.Repeat(", ?", len(productIDs)-1) // Products ID args placeholder.
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// Convert productIDs to []interface{} (any interface)
	// Creating a list of product IDs.
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	// Spread all product IDs as arguments in SELECT Query.
	// replacing the above mentioned `placeholder` "?"s.
	rows, err := s.db.Query(query, args...) // Query executed. Rows returned.
	if err != nil {
		return nil, err
	}

	products := []types.Product{} // List to store list(slice) of `types.Product`.

	// Iterating over rows from query.
	for rows.Next() {
		p, err := scanRowIntoProduct(rows) // transforming each row into `types.Product`.
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

// Register a new product.
// TODO : Create endpoint with JWT auth middleware. Only `admin` access.
func (s *Store) RegisterProduct(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	return nil
}

// Create a `types.Product` struct from Query returned rows...
// Only scans first row.
func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return product, nil
}

// Update product values in DB.
func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, image = ?, description = ?, quantity = ? WHERE id = ?", product.Name, product.Image, product.Description, product.Quantity, product.ID)

	if err != nil {
		return err
	}

	return nil
}
