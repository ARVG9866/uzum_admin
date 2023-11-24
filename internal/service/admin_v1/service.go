package service

import (
	"context"

	"github.com/Shemistan/uzum_admin/internal/models"
	"github.com/Shemistan/uzum_admin/internal/storage"
)

type IService interface {
	AddProduct(context.Context, *models.Product) error
	UpdateProduct(context.Context, *models.Product) error
	DeleteProduct(context.Context, int64) error
	GetProduct(context.Context, int64) (*models.Product, error)
	GetAllProducts(context.Context, *models.GetAllSetting) ([]*models.Product, error)
	GetStatistics(context.Context) (*models.Statistic, error)
}

func NewService(storage storage.IStorage) IService {
	return &service{
		storage: storage,
	}
}

type service struct {
	storage storage.IStorage
}

func (s *service) AddProduct(ctx context.Context, product *models.Product) error {
	if err := s.storage.CreateProduct(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateProduct(ctx context.Context, product *models.Product) error {
	if err := s.storage.UpdateProduct(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteProduct(ctx context.Context, product_id int64) error {
	if err := s.storage.DeleteProduct(ctx, product_id); err != nil {
		return err
	}

	return nil
}
func (s *service) GetProduct(ctx context.Context, product_id int64) (*models.Product, error) {
	product, err := s.storage.GetProduct(ctx, product_id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) GetAllProducts(ctx context.Context, setting *models.GetAllSetting) ([]*models.Product, error) {
	products, err := s.storage.GetAllProducts(ctx, setting)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *service) GetStatistics(ctx context.Context) (*models.Statistic, error) {
	statistic, err := s.storage.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return statistic, nil
}
