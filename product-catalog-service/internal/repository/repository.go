package repository

import (
	"context"
	"database/sql"
	"product-catalog-service/internal/entity"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	var product *entity.Product

	query := `SELECT id, name, description, price, stock FROM products WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	query := `INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Stock)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.ID = int(id)
	return product, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	query := `UPDATE products SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Stock, product.ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetProducts(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product

	query := `SELECT id, name, description, price, stock FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
