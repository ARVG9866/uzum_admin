package storage

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"

	"github.com/Shemistan/uzum_admin/internal/models"
)

const productTable = "product"
const statisticTable = "statistic"
const orderProductTable = "order_product"

type IStorage interface {
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
	m := structs.Map(req)
	builder := sq.Update(productTable).SetMap(m).
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
	builder := sq.Update(productTable).SetMap(map[string]interface{}{
		"count":   0,
		"deleted": true,
	}).
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
		Where(sq.Eq{"id": id, "deleted": false}).
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
		Where(sq.Eq{"deleted": false}).
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
	var statistic models.Statistic

	str_query := fmt.Sprintf(`SELECT SUM(count), SUM(price) FROM %s;`, orderProductTable)

	err := s.db.QueryRowContext(ctx, str_query).Scan(&statistic.CountSold, &statistic.Earned)
	if err != nil {
		return nil, err
	}

	str_query = fmt.Sprintf(`SELECT %s.id, %s.name, %s.price, %s.count FROM %s
			JOIN %s ON %s.product_id=%s.id;`,
		productTable, productTable, orderProductTable, orderProductTable, orderProductTable, productTable, orderProductTable, productTable)

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

	return &statistic, nil
}
