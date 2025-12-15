package event

import (
	"strings"
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

func HandleShearCommand(s *discordgo.Session, message *discordgo.MessageCreate, commandSuffix string, executor CommandHandler) {
	// ignore all messages created by the bot itself
	if message.Author.ID == s.State.User.ID {
		return
	}

	fullCommandPrefix := CommandPrefix + " " + commandSuffix
	if strings.HasPrefix(strings.ToLower(message.Content), strings.ToLower(fullCommandPrefix)) {
		logger.Info("Executing message command", "command", commandSuffix, "user", message.Author.Username)
		arg := strings.TrimSpace(message.Content[len(fullCommandPrefix):])

		executor(s, message.Message, arg)
	}
}

func findMemberIDByUsername(s *discordgo.Session, guildID, username string) (string, error) {
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		logger.Error("Failed to fetch guild members", "Error", err)
	}

	for _, member := range members {

		// check primary username (member.User.Username)
		if strings.EqualFold(member.User.Username, username) {
			return member.User.ID, nil
		}

		// heck guild nickname (member.Nick)
		if member.Nick != "" && strings.EqualFold(member.Nick, username) {
			return member.User.ID, nil
		}
	}

	return "", err
}
