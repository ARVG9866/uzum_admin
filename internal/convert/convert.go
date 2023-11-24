package convert

import (
	"github.com/Shemistan/uzum_admin/internal/models"
	pb "github.com/Shemistan/uzum_admin/pkg/admin_v1"
)

func PbToModelProduct(product *pb.Product) *models.Product {
	return &models.Product{
		ID:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Count:       product.Count,
	}
}

func ModelToPbProduct(product *models.Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Count:       product.Count,
	}
}

func ModelToPbStatistic(statistic *models.Statistic, products []*pb.Product) *pb.Statistic {
	return &pb.Statistic{
		CountSold: statistic.CountSold,
		Earned:    statistic.Earned,
		Product:   products,
	}
}
