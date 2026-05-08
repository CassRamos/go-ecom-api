-- name: ListProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products(name, description, price_in_cents, quantity) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price_in_cents = $4, quantity = $5
WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: DecrementProductQuantity :exec
UPDATE products
SET quantity = quantity - $2
WHERE id = $1;

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrderItemsByOrderId :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) 
VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents)
VALUES ($1, $2, $3, $4) RETURNING *;