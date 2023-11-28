package api

import (
	"context"

	"github.com/Shemistan/uzum_admin/internal/convert"
	"github.com/Shemistan/uzum_admin/internal/models"
	admin_v1 "github.com/Shemistan/uzum_admin/internal/service/admin_v1"
	pb "github.com/Shemistan/uzum_admin/pkg/admin_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Admin struct {
	pb.UnimplementedAdminV1Server

	AdminService admin_v1.IService
}

func (a *Admin) AddProduct(ctx context.Context, req *pb.AddProduct_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = a.AdminService.AddProduct(ctx, convert.PbToModelProduct(req.Product))

	return &emptypb.Empty{}, err
}

func (a *Admin) UpdateProduct(ctx context.Context, req *pb.UpdateProduct_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}
	err = a.AdminService.UpdateProduct(ctx, convert.PbToModelProduct(req.Product))

	return &emptypb.Empty{}, err
}

func (a *Admin) DeleteProduct(ctx context.Context, req *pb.DeleteProduct_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = a.AdminService.DeleteProduct(ctx, req.ProductId)

	return &emptypb.Empty{}, err
}

func (a *Admin) GetProduct(ctx context.Context, req *pb.GetProduct_Request) (*pb.GetProduct_Response, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	product, err := a.AdminService.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.GetProduct_Response{
		Product: convert.ModelToPbProduct(product),
	}, err
}

func (a *Admin) GetProducts(ctx context.Context, req *pb.GetProducts_Request) (*pb.GetProducts_Response, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	setting := &models.GetAllSetting{
		Page: req.Page,
		Size: req.Size,
	}

	res, err := a.AdminService.GetAllProducts(ctx, setting)
	if err != nil {
		return nil, err
	}

	products := make([]*pb.Product, 0, len(res))

	for _, v := range res {
		products = append(products, convert.ModelToPbProduct(v))
	}

	rtn := &pb.GetProducts_Response{
		Product: products,
	}

	return rtn, nil
}

func (a *Admin) GetStatistics(ctx context.Context, _ *emptypb.Empty) (*pb.GetStatistics_Response, error) {
	res, err := a.AdminService.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*pb.Product, 0, len(res.Products))

	for _, v := range res.Products {
		products = append(products, convert.ModelToPbProduct(v))
	}

	rtn := &pb.GetStatistics_Response{
		Statistic: convert.ModelToPbStatistic(res, products),
	}

	return rtn, nil
}
