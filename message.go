package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	// filter out bot messages
	if message.Author.Bot {
		return
	}

	updatedTime := time.Now().Format("2006-01-02")
	logger.Info("User sent a message", "User", message.Author.Username, "Time", updatedTime)

	go WriteToPostgres(conn, ctx, "message", updatedTime, message.Author.Username)
}
