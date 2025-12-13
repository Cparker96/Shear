package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
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

	go WriteToPostgres(conn, ctx, "message", updatedTime, message.Author.Username)
}
