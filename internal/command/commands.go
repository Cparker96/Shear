package command

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, message *discordgo.Message, arg string)

func executeGetUserActivity(s *discordgo.Session, message *discordgo.Message, arg string) {
	ctx := context.Background()
	pool := conn

	// send initial status message
	thinkingMsg, err := s.ChannelMessageSend(message.ChannelID, "Fetching activity and preparing CSV file...")
	if err != nil {
		logger.Error("Failed to send thinking message", "error", err)
		return
	}

	query := "SELECT * FROM public.activity"
	rows, err := pool.Query(ctx, query)
	if err != nil {
		logger.Error("Failed to query for recent activity", "Error", err)
		s.ChannelMessageEdit(message.ChannelID, thinkingMsg.ID, "Database Error! Could not execute query")
		return
	}
	defer rows.Close()

	var bytes bytes.Buffer
	writer := csv.NewWriter(&bytes)

	columns := rows.FieldDescriptions()
	headers := make([]string, len(columns))
	for index, fd := range columns {
		headers[index] = string(fd.Name)
	}

	if err := writer.Write(headers); err != nil {
		logger.Error("Error writing CSV header", "Error", err)
		s.ChannelMessageEdit(message.ChannelID, thinkingMsg.ID, "Internal Error during file preparation (Header).")
		return
	}

	count := 0
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			logger.Error("Error getting row values", "error", err)
			continue
		}

		// convert all values to strings for CSV writing
		record := make([]string, len(values))
		for index, value := range values {
			if time, ok := value.(time.Time); ok {
				record[index] = time.Format("2006-01-02")
			} else if value != nil {
				record[index] = fmt.Sprintf("%v", value)
			} else {
				record[index] = "" // handle nulls
			}
		}

		if err := writer.Write(record); err != nil {
			logger.Error("Error writing CSV record", "Error", err)
			continue
		}
		count++
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		logger.Error("Error flushing CSV writer", "Error", err)
		s.ChannelMessageEdit(message.ChannelID, thinkingMsg.ID, "Internal Error during file finalization.")
		return
	}

	// upload message to discord
	if count == 0 {
		s.ChannelMessageEdit(message.ChannelID, thinkingMsg.ID, "Query successful, but no activities found in the table.")
		return
	}

	fileName := fmt.Sprintf("activity_log_%s.csv", time.Now().Format("20060102_150405"))

	// use ChannelFileSend to upload the buffer contents
	_, err = s.ChannelFileSend(message.ChannelID, fileName, &bytes)
	if err != nil {
		logger.Error("Failed to send CSV file to Discord", "error", err)
		s.ChannelMessageEdit(message.ChannelID, thinkingMsg.ID, fmt.Sprintf("File upload failed after generating %d rows.", count))
		return
	}

	// final success message (delete the thinking message, or edit it)
	s.ChannelMessageEdit(message.ChannelID, thinkingMsg.ID, fmt.Sprintf("Attached CSV file containing **%d** activity records.", count))
}

func executeRemoveUser(s *discordgo.Session, message *discordgo.Message, arg string) {
	username := arg
	ctx := context.Background()
	pool := conn

	query := fmt.Sprintf("DELETE FROM public.activity where username = '%s'", username)

	commandTag, err := pool.Exec(ctx, query)
	if err != nil {
		logger.Error("Error executing query", "error", err)
		return
	}

	if commandTag.RowsAffected() == 0 {
		logger.Warn("No rows updated", "userID", username)
	}

	logger.Info("User activity updated", "Username", username, "Action", "RemovedUser", "rowsAffected", commandTag.RowsAffected())

	memberID, err := findMemberIDByUsername(s, message.GuildID, username)
	if err != nil {
		logger.Warn("Failed to find member ID for kick attempt", "Username", username, "Error", err)
		return
	}

	kickErr := s.GuildMemberDelete(message.GuildID, memberID)
	if kickErr != nil {
		logger.Error("Failed to kick user from Discord", "UserID", memberID, "Username", username, "Error", kickErr)
		return
	}

	s.ChannelMessageSend(message.ChannelID, "Successfully removed user from DB and kicked from Discord")
	logger.Info("Successfully removed user from DB and kicked from Discord", "Username", username, "Invoker", message.Author.Username)
}
