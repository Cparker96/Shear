package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MessageReaction(s *discordgo.Session, reaction *discordgo.MessageReactionAdd, conn *pgxpool.Pool, logger *slog.Logger) {
	ctx := context.Background()
	user, err := s.User(reaction.UserID)
	if err != nil {
		logger.Error("Error getting user", "Error", err)
		return
	}

	// filter out bot reactions
	if user.Bot {
		return
	}

	updatedTime := time.Now().Format("2006-01-02")
	logger.Info("User reacted to a message", "User", user.Username, "Time", updatedTime)

	go database.WriteToPostgres(conn, ctx, "reaction", updatedTime, user.Username, logger)
}
