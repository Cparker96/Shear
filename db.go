package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserEvent struct {
	Username   string `json:"username"`
	UpdateType string `json:"update_type"`
	Date       string `json:"date"`
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

func WriteToPostgres(pool *pgxpool.Pool, ctx context.Context, action string, date string, user string, logger *slog.Logger) error {
	checkForUser, err := DoesUserExist(pool, ctx, user)
	if err != nil {
		logger.Error("Failed to retrieve query results", "Error", err)
		return err
	}

	if checkForUser.Username == user {
		query := `
			UPDATE public.activity
			SET update_type = $1, date = $2
			WHERE username = $3
		`
		// update existing record
		UpsertUserActivity(pool, ctx, query, action, date, user)
	} else {
		query := `
			INSERT INTO public.activity
			VALUES ($3, $1, $2)
		`
		// insert new record
		UpsertUserActivity(pool, ctx, query, action, date, user)
	}

	return nil
}

func DoesUserExist(pool *pgxpool.Pool, ctx context.Context, user string) (UserEvent, error) {
	query := `
        SELECT username, update_type, date 
        FROM public.activity
        WHERE username = $1
		LIMIT 1
    `

	rows, err := pool.Query(ctx, query, user)
	if err != nil {
		logger.Error("Error executing SELECT query", "error", err)
		return UserEvent{}, err
	}
	defer rows.Close()

	event := UserEvent{}
	if rows.Next() {
		err := rows.Scan(&event.Username, &event.UpdateType, &event.Date)
		if err != nil {
			logger.Error("Failed to retrieve user event", "Username", user)
			return UserEvent{}, err
		}
		return event, nil
	}

	return event, nil
}

func UpsertUserActivity(pool *pgxpool.Pool, ctx context.Context, query string, action string, date string, user string) {
	commandTag, err := pool.Exec(ctx, query, action, date, user)
	if err != nil {
		logger.Error("Error executing query", "error", err)
	}

	if commandTag.RowsAffected() == 0 {
		logger.Warn("No rows updated", "userID", user)
	}

	logger.Info("User activity updated", "Username", user, "Action", action, "rowsAffected", commandTag.RowsAffected())
}
