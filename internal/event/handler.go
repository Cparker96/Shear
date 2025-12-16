package event

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandHandler func(s *discordgo.Session, message *discordgo.Message, conn *pgxpool.Pool, arg string, logger *slog.Logger)

func HandleShearCommand(s *discordgo.Session, message *discordgo.MessageCreate, conn *pgxpool.Pool, CommandPrefix string, CommandSuffix string, executor CommandHandler, logger *slog.Logger) {
	// ignore all messages created by the bot itself
	if message.Author.ID == s.State.User.ID {
		return
	}

	fullCommandPrefix := CommandPrefix + " " + CommandSuffix
	if strings.HasPrefix(strings.ToLower(message.Content), strings.ToLower(fullCommandPrefix)) {
		logger.Info("Executing message command", "command", CommandSuffix, "user", message.Author.Username)
		arg := strings.TrimSpace(message.Content[len(fullCommandPrefix):])

		executor(s, message.Message, conn, arg, logger)
	}
}

func FindMemberIDByUsername(s *discordgo.Session, guildID, username string, logger *slog.Logger) (string, error) {
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		logger.Error("Failed to fetch guild members", "Error", err)
	}

	for _, member := range members {

		// check primary username (member.User.Username)
		if strings.EqualFold(member.User.Username, username) {
			return member.User.ID, nil
		}

		// check guild nickname (member.Nick)
		if member.Nick != "" && strings.EqualFold(member.Nick, username) {
			return member.User.ID, nil
		}
	}

	return "", err
}
