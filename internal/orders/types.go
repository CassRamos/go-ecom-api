package orders

import (
	"context"

	repository "github.com/CassRamos/go-ecom-api.git/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customer_id"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repository.Order, error)
}
