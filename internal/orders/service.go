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

	// Create a new query object bound to the transaction, bypassing the interface limitation
	qtx := repository.New(tx)

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

		err = qtx.DecrementProductQuantity(ctx, repository.DecrementProductQuantityParams{
			ID:       item.ProductId,
			Quantity: item.Quantity,
		})
		if err != nil {
			return repository.Order{}, fmt.Errorf("failed to decrement stock: %w", err)
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

func (s *svc) CancelOrder(ctx context.Context, orderId int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	qtx := repository.New(tx)

	order, err := qtx.GetOrderById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("Order not found: %w", err)
	}

	if order.Status == "CANCELED" {
		return fmt.Errorf("order is already canceled")
	}

	err = qtx.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
		ID:     order.ID,
		Status: "CANCELED",
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	items, err := qtx.GetOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}

	for _, item := range items {
		err = qtx.IncrementProductQuantity(ctx, repository.IncrementProductQuantityParams{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
		if err != nil {
			return fmt.Errorf("failed to restore stock for product %d: %w", item.ProductID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *svc) GetOrderByID(ctx context.Context, id int64) (OrderResponse, error) {
	order, err := s.repo.GetOrderById(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	items, err := s.repo.GetOrderItemsByOrderId(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	return OrderResponse{
		Order: order,
		Items: items,
	}, nil
}
