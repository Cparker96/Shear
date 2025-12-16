package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MessageCreate(s *discordgo.Session, message *discordgo.MessageCreate, conn *pgxpool.Pool, logger *slog.Logger) {
	ctx := context.Background()
	channel, err := s.Channel(message.ChannelID)
	if err != nil {
		logger.Error("Error getting channel", "Error", err)
		return
	}

	// only track regular text channels
	if channel.Type != discordgo.ChannelTypeGuildText {
		logger.Debug("Ignoring non-text channel", "Type", channel.Type)
		return
	}

	// filter out bot messages
	if message.Author.Bot {
		return
	}

	// filter out webhook messages (followed channels, integrations)
	if message.WebhookID != "" {
		return
	}

	updatedTime := time.Now().Format("2006-01-02")
	logger.Info("User sent a message", "User", message.Author.Username, "Time", updatedTime)

	go database.WriteToPostgres(conn, ctx, "message", updatedTime, message.Author.Username, logger)
}
