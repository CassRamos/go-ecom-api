package orders

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/CassRamos/go-ecom-api.git/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock for product")
)

type svc struct {
	repo *repository.Queries
	db   *pgx.Conn
}

func NewService(repo *repository.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repository.Order, error) {
	if tempOrder.CustomerID == 0 {
		return repository.Order{}, fmt.Errorf("customer ID is required")
	}

	if len(tempOrder.Items) == 0 {
		return repository.Order{}, fmt.Errorf("at least one order item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repository.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repository.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range tempOrder.Items {
		product, err := qtx.GetProductByID(ctx, item.ProductId)
		if err != nil {
			return repository.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repository.Order{}, ErrInsufficientStock
		}

		_, err = qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductId,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return repository.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	tx.Commit(ctx)

	return order, nil
}
