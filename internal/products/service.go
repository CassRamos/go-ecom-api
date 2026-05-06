package products

import (
	"context"

	repository "github.com/CassRamos/go-ecom-api.git/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repository.Product, error)
	CreateProduct(ctx context.Context, params repository.CreateProductParams) (repository.Product, error)
	UpdateProduct(ctx context.Context, params repository.UpdateProductParams) (repository.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type svc struct {
	repo repository.Querier
}

func NewService(repo repository.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListProducts(ctx context.Context) ([]repository.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) CreateProduct(ctx context.Context, params repository.CreateProductParams) (repository.Product, error) {
	return s.repo.CreateProduct(ctx, params)
}

func (s *svc) UpdateProduct(ctx context.Context, params repository.UpdateProductParams) (repository.Product, error) {
	return s.repo.UpdateProduct(ctx, params)
}

func (s *svc) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}
