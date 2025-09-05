package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/zippy/internal/entity"
)

type Product struct {
	DB *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) *Product {
	return &Product{
		DB: db,
	}
}

func (p *Product) Create(product *entity.Product) error {
	_, err := p.DB.Exec(context.Background(),
		"INSERT INTO products (id, name, price, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)",
		product.Id,
		product.Name,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (p *Product) FetchPaged(pageNumber, pageSize int, sort string) ([]entity.Product, error) {
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	products := make([]entity.Product, 0)
	offset := (pageNumber - 1) * pageSize
	rows, err := p.DB.Query(context.Background(), fmt.Sprintf("SELECT id, name, price, created_at, updated_at FROM products ORDER BY name %v LIMIT $1 OFFSET $2", sort), pageSize, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *Product) FetchById(id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	rows := p.DB.QueryRow(context.Background(), "SELECT id, name, price, created_at, updated_at FROM products WHERE id = $1", id)
	err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt)

	return &product, err
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.FetchById(product.Id)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(context.Background(),
		"UPDATE products SET (name, price, updated_at) = ($1, $2, $3) WHERE id = $4",
		product.Name,
		product.Price,
		product.UpdatedAt,
		product.Id,
	)
	return err
}

func (p *Product) Delete(id uuid.UUID) error {
	product, err := p.FetchById(id)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(context.Background(), "DELETE FROM products WHERE id = $1", product.Id)
	return err
}
