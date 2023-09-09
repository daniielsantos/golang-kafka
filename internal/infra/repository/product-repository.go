package repository

import (
	"database/sql"
	"fmt"

	"github.com/daniielsantos/dss/internal/entity"
)

type ProductRepository struct {
	Db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{Db: db}
}

func (r *ProductRepository) Create(product *entity.Product) error {
	_, err := r.Db.Exec("INSERT INTO product (id, name, price) VALUES (?,?,?)",
		product.Id, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) FindAll() ([]*entity.Product, error) {
	rows, err := r.Db.Query("SELECT id, name, price FROM product")
	if err != nil {
		fmt.Printf("bosta ", err)
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
