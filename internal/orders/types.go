package orders

import (
	"context"

	repository "github.com/CassRamos/go-ecom-api.git/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductId int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repository.Order, error)
	GetOrderByID(ctx context.Context, id int64) (OrderResponse, error)
}

type OrderResponse struct {
	Order repository.Order       `json:"order"`
	Items []repository.OrderItem `json:"items"`
}
