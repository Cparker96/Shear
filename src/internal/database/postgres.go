package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserEvent struct {
	Username   string `json:"username"`
	UpdateType string `json:"update_type"`
	Date       string `json:"date"`
}

func PostgresConn(postgresUser string, postgresPW string, postgresDBName string) (*pgxpool.Pool, context.Context, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s", postgresUser, postgresPW, postgresDBName)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect to postgres", "Error", err)
	}

	return pool, ctx, nil
}

func WriteToPostgres(pool *pgxpool.Pool, ctx context.Context, action string, date string, user string) error {
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
		return
	}

	if commandTag.RowsAffected() == 0 {
		logger.Warn("No rows updated", "userID", user)
	}

	logger.Info("User activity updated", "Username", user, "Action", action, "rowsAffected", commandTag.RowsAffected())
}
