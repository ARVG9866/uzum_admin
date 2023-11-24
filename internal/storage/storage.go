package storage

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/Shemistan/uzum_admin/internal/models"
)

const productTable = "product"
const statisticTable = "statistic"
const historyTable = "history"

type IStorage interface {
	// AddProduct(ctx context.Context, req *models.Product) (int64, error)
	CreateProduct(ctx context.Context, req *models.Product) error
	UpdateProduct(ctx context.Context, req *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context, req *models.GetAllSetting) ([]*models.Product, error)
	GetStatistics(ctx context.Context) (*models.Statistic, error)
}

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) IStorage {
	return &storage{db: db}
}

func (s *storage) CreateProduct(ctx context.Context, req *models.Product) error {
	builder := sq.Insert(productTable).
		Columns("name", "description", "price", "count").
		Values(req.Name, req.Description, req.Price, req.Count).
		RunWith(s.db).
		PlaceholderFormat(sq.Dollar)

	_, err := builder.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) UpdateProduct(ctx context.Context, req *models.Product) error {
	builder := sq.Update(productTable).SetMap(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"count":       req.Count,
	}).
		Where(sq.Eq{"id": req.ID}).
		RunWith(s.db).
		PlaceholderFormat(sq.Dollar)

	_, err := builder.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteProduct(ctx context.Context, id int64) error {
	builder := sq.Delete(productTable).
		Where(sq.Eq{"id": id}).
		RunWith(s.db).
		PlaceholderFormat(sq.Dollar)

	_, err := builder.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product

	builder := sq.Select("id", "name", "description", "price", "count").
		From(productTable).
		Where(sq.Eq{"id": id}).
		RunWith(s.db).
		PlaceholderFormat(sq.Dollar)

	err := builder.QueryRowContext(ctx).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Count)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *storage) GetAllProducts(ctx context.Context, req *models.GetAllSetting) ([]*models.Product, error) {
	var products []*models.Product

	builder := sq.Select("id", "name", "description", "price", "count").
		From(productTable).
		RunWith(s.db).
		PlaceholderFormat(sq.Dollar).
		Offset(uint64(req.Page)).
		Limit(uint64(req.Size))

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product

		if err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Count); err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (s *storage) GetStatistics(ctx context.Context) (*models.Statistic, error) {
	var statistic *models.Statistic

	str_query := fmt.Sprintf(`SELECT SUM(count), SUM(price) FROM %s`, historyTable)

	err := s.db.QueryRowContext(ctx, str_query).Scan(&statistic.CountSold, &statistic.Earned)
	if err != nil {
		return nil, err
	}

	str_query = fmt.Sprintf(`SELECT %s.id, %s.name, %s.price, %s.count FROM %s
			JOIN %s USING (product_ID)`,
		productTable, productTable, historyTable, historyTable, historyTable, productTable)

	rows, err := s.db.QueryContext(ctx, str_query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product

		if err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Count); err != nil {
			return nil, err
		}

		statistic.Products = append(statistic.Products, &product)
	}

	return statistic, nil
}
