package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/CassRamos/go-ecom-api.git/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Create a context with a timeout of 5 seconds. This context will be used for the database connection attempt. If the connection takes longer than 5 seconds, it will be canceled.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := env.GetString("DATABASE_URL", "")
	if dsn == "" {
		panic("DATABASE_URL is not set")
	}

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: dsn,
		},
	}

	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database connection
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}

	defer conn.Close(ctx)

	logger.Info("Connected to database")

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
