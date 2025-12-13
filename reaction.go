// reaction.go - Reaction tracking logic
package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageReaction(s *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
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

	go WriteToPostgres(conn, ctx, "reaction", updatedTime, user.Username)
}
