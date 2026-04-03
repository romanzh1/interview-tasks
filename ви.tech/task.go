package main

import (
	"database/sql"
	"errors"
)

type ProductService struct {
	db          *sql.DB
	orderClient *OrderClient
}

type Product struct {
	ID       int64   // id
	Name     string  // наименование
	Price    float64 // цена
	Quantity int     // количество
}

func (s *ProductService) CreateOrder(product_id int64) error {
	var product Product

	row := s.db.QueryRow(
		"SELECT * FROM products WHERE id = ?",
		product_id,
	)

	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
		return err
	}

	if product.Quantity == 0 {
		return errors.New("product is out of stock")
	}

	_, err := s.db.Exec(
		"UPDATE products SET quantity = $1 WHERE id = $2",
		product.Quantity-1,
		product_id,
	)
	if err != nil {
		return err
	}

	return s.orderClient.CreateOrder(product.ID)
}

// заглушка, так как на скрине нет реализации
type OrderClient struct{}

func (c *OrderClient) CreateOrder(productID int64) error {
	return nil
}
