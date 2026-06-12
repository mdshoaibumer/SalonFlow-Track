package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// ProductService exposes product/inventory operations to the Wails frontend.
type ProductService struct {
	ctx      context.Context
	uc       *usecase.ProductUseCase
	guard    *PermissionGuard
	licGuard *LicenseGuard
}

func NewProductService(uc *usecase.ProductUseCase) *ProductService {
	return &ProductService{uc: uc}
}

func (s *ProductService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ProductService) ListProducts(input usecase.ListProductsInput) (*usecase.ListProductsOutput, error) {
	return s.uc.ListProducts(s.ctx, input)
}

func (s *ProductService) GetProduct(id string) (*domain.Product, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetProduct(s.ctx, uid)
}

func (s *ProductService) CreateProduct(input usecase.CreateProductInput) (*domain.Product, error) {
	return s.uc.CreateProduct(s.ctx, input)
}

func (s *ProductService) UpdateProduct(id string, input usecase.UpdateProductInput) (*domain.Product, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.UpdateProduct(s.ctx, uid, input)
}

func (s *ProductService) DeleteProduct(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteProduct(s.ctx, uid)
}

func (s *ProductService) AdjustStock(input usecase.StockAdjustInput) (*domain.StockTransaction, error) {
	if err := s.licGuard.RequireActive(domain.OpInventoryChange); err != nil {
		return nil, err
	}
	return s.uc.AdjustStock(s.ctx, input)
}

func (s *ProductService) ListStockHistory(input usecase.ListStockHistoryInput) (*usecase.ListStockHistoryOutput, error) {
	return s.uc.ListStockHistory(s.ctx, input)
}

func (s *ProductService) CreatePurchase(input usecase.CreatePurchaseInput) (*domain.PurchaseEntry, error) {
	return s.uc.CreatePurchase(s.ctx, input)
}

func (s *ProductService) GetPurchase(id string) (*domain.PurchaseEntry, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetPurchase(s.ctx, uid)
}

func (s *ProductService) ListPurchases(input usecase.ListPurchasesInput) (*usecase.ListPurchasesOutput, error) {
	return s.uc.ListPurchases(s.ctx, input)
}

func (s *ProductService) GetInventoryStats() (*domain.InventoryStats, error) {
	return s.uc.GetInventoryStats(s.ctx)
}

func (s *ProductService) GetLowStockProducts() ([]domain.LowStockItem, error) {
	return s.uc.GetLowStockProducts(s.ctx)
}
