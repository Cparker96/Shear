package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityData struct {
	Username   string `json:"username"`
	UpdateType string `json:"update_type"`
}

func PostgresConn(postgresPW string, logger *slog.Logger) (*pgxpool.Pool, context.Context, error) {
	dsn := fmt.Sprintf("postgresql://shear_user:%s@localhost:5432/shear", postgresPW)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect to postgres", "Error", err)
	}

	return pool, ctx, nil
}

func FindData(pool *pgxpool.Pool, ctx context.Context, logger *slog.Logger) (ActivityData, error) {
	query := fmt.Sprintf("SELECT * FROM public.activity")
	rows, err := pool.Query(ctx, query)
	if err != nil {
		logger.Error("Error executing query", "Error", err)
	}
	defer rows.Close()

	data := ActivityData{}
	if rows.Next() {
		err := rows.Scan(&data.Username, &data.UpdateType)
		if err != nil {
			logger.Error("Failed to find data", "Error", err)
		}
		return data, nil
	}

	return data, nil
}
