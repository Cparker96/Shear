package database

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

func PostgresConn(postgresUser string, postgresPW string, postgresDBName string, logger *slog.Logger) (*pgxpool.Pool, context.Context, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s", postgresUser, postgresPW, postgresDBName)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect to postgres", "Error", err)
	}

	return pool, ctx, nil
}

func WriteToPostgres(pool *pgxpool.Pool, ctx context.Context, action string, date string, user string, logger *slog.Logger) error {
	checkForUser, err := DoesUserExist(pool, ctx, user, logger)
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
		UpsertUserActivity(pool, ctx, query, action, date, user, logger)
	} else {
		query := `
			INSERT INTO public.activity (username, update_type, date)
			VALUES ($3, $1, $2)
		`
		// insert new record
		UpsertUserActivity(pool, ctx, query, action, date, user, logger)
	}

	return nil
}

func DoesUserExist(pool *pgxpool.Pool, ctx context.Context, user string, logger *slog.Logger) (UserEvent, error) {
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

func UpsertUserActivity(pool *pgxpool.Pool, ctx context.Context, query string, action string, date string, user string, logger *slog.Logger) {
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

func GetAllUserEvents(pool *pgxpool.Pool, ctx context.Context, logger *slog.Logger) ([]UserEvent, error) {
	query := `
		SELECT username, update_type, date
		FROM public.activity
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		logger.Error("Error executing SELECT query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var events []UserEvent
	for rows.Next() {
		event := UserEvent{}
		err := rows.Scan(&event.Username, &event.UpdateType, &event.Date)
		if err != nil {
			logger.Error("Failed to retrieve user event", "error", err)
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func HasRecords(pool *pgxpool.Pool, ctx context.Context, logger *slog.Logger) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM public.activity
	`

	var count int
	err := pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		logger.Error("Error checking for existing records", "error", err)
		return false, err
	}

	return count > 0, nil
}

func SeedDatabase(pool *pgxpool.Pool, ctx context.Context, usernames []string, date string, logger *slog.Logger) error {
	if len(usernames) == 0 {
		logger.Info("No usernames provided for seeding")
		return nil
	}

	query := `
		INSERT INTO public.activity (username, update_type, date)
		VALUES ($1, 'seed', $2)
	`

	inserted := 0
	skipped := 0
	for _, username := range usernames {
		// Check if user already exists before inserting
		existingUser, err := DoesUserExist(pool, ctx, username, logger)
		if err != nil {
			logger.Warn("Failed to check if user exists during seeding", "username", username, "error", err)
			continue
		}

		if existingUser.Username == username {
			skipped++
			continue
		}

		commandTag, err := pool.Exec(ctx, query, username, date)
		if err != nil {
			logger.Warn("Failed to insert username during seeding", "username", username, "error", err)
			continue
		}
		if commandTag.RowsAffected() > 0 {
			inserted++
		}
	}

	logger.Info("Database seeding completed", "total_usernames", len(usernames), "inserted", inserted, "skipped", skipped, "seed_date", date)
	return nil
}
