package event

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DetectUserJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd, conn *pgxpool.Pool, logger *slog.Logger) {
	ctx := context.Background()
	logger.Info("New member joined", "Username", m.User.Username)

	updatedTime := time.Now().Format("2006-01-02")
	go database.WriteToPostgres(conn, ctx, "join", updatedTime, m.User.Username, logger)
}

func DetectUserLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove, conn *pgxpool.Pool, logger *slog.Logger) {
	ctx := context.Background()
	pool := conn
	logger.Info("Member left", "Username", m.User.Username)

	query := fmt.Sprintf("DELETE FROM public.activity where username = '%s'", m.User.Username)

	commandTag, err := pool.Exec(ctx, query)
	if err != nil {
		logger.Error("Error executing query", "Error", err)
		return
	}

	if commandTag.RowsAffected() == 0 {
		logger.Warn("No rows updated", "Username", m.User.Username)
	}

	logger.Info("User activity updated", "Username", m.User.Username, "Action", "RemovedUser", "rowsAffected", commandTag.RowsAffected())
}
