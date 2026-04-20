package scheduler

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduledCronJob struct {
	conn      *pgxpool.Pool
	logger    *slog.Logger
	session   *discordgo.Session
	channelID string
	roleName  string
}

type InactiveUser struct {
	Username   string
	UpdateType string
	UserDate   string
	DaysDelta  int
}

func NewScheduledCronJob(conn *pgxpool.Pool, logger *slog.Logger, session *discordgo.Session, channelID string, roleName string) *ScheduledCronJob {
	return &ScheduledCronJob{
		conn:      conn,
		logger:    logger,
		session:   session,
		channelID: channelID,
		roleName:  roleName,
	}
}

func (cron *ScheduledCronJob) GetUserActivityWithTimeRanges() {
	cron.logger.Info("Initializing scheduled cron job for fetching user activity time deltas in Discord")

	ctx := context.Background()

	events, err := database.GetAllUserEvents(cron.conn, ctx, cron.logger)
	if err != nil {
		cron.logger.Error("Failed to fetch user events", "error", err)
		return
	}

	today := time.Now()
	todayStr := today.Format("2006-01-02")
	todayDate, err := time.Parse("2006-01-02", todayStr)
	if err != nil {
		cron.logger.Error("Failed to parse today's date", "error", err)
		return
	}

	var inactiveUsers []InactiveUser

	for _, event := range events {
		if event.Date == "" {
			cron.logger.Warn("User has no date set, skipping inactivity check", "username", event.Username)
			continue
		}

		userDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			cron.logger.Warn("Failed to parse user date", "username", event.Username, "date", event.Date, "error", err)
			continue
		}

		delta := todayDate.Sub(userDate)
		daysDelta := int(delta.Hours() / 24)

		cron.logger.Info("User activity time delta calculated",
			"username", event.Username,
			"update_type", event.UpdateType,
			"user_date", event.Date,
			"today", todayStr,
			"days_delta", daysDelta,
		)

		if daysDelta > 120 {
			inactiveUsers = append(inactiveUsers, InactiveUser{
				Username:   event.Username,
				UpdateType: event.UpdateType,
				UserDate:   event.Date,
				DaysDelta:  daysDelta,
			})
		}
	}

	if len(inactiveUsers) > 0 {
		if cron.channelID == "" {
			cron.logger.Warn("Inactive users found but no channel ID configured", "count", len(inactiveUsers))
			return
		}

		channel, err := cron.session.Channel(cron.channelID)
		if err != nil {
			cron.logger.Error("Failed to get channel", "error", err)
			return
		}

		roleMention := cron.findRoleMention(channel.GuildID, cron.roleName)

		// Build CSV
		csvBuffer, err := cron.buildCSV(inactiveUsers)
		if err != nil {
			cron.logger.Error("Failed to build CSV", "error", err)
			return
		}

		// Send summary message
		summary := fmt.Sprintf(
			"Users exceeding 120 days of inactivity: **%d**\nReport Date: %s\nSee attached CSV.",
			len(inactiveUsers),
			todayStr,
		)

		if roleMention != "" {
			summary = roleMention + "\n\n" + summary
		}

		_, err = cron.session.ChannelMessageSend(cron.channelID, summary)
		if err != nil {
			cron.logger.Error("Failed to send summary message", "error", err)
			return
		}

		// Send CSV file
		_, err = cron.session.ChannelFileSend(
			cron.channelID,
			"inactive_users.csv",
			csvBuffer,
		)
		if err != nil {
			cron.logger.Error("Failed to send CSV file", "error", err)
			return
		}

		cron.logger.Info("Sent inactive users CSV to Discord", "count", len(inactiveUsers))
	} else {
		cron.logger.Info("No inactive users found exceeding 120 day threshold")
	}

	cron.logger.Info("Completed scheduled cron job for user activity time deltas",
		"total_records", len(events),
		"inactive_users", len(inactiveUsers),
	)
}

func (cron *ScheduledCronJob) buildCSV(users []InactiveUser) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	if err := writer.Write([]string{"Username", "Last Activity", "Update Type", "Days Inactive"}); err != nil {
		return nil, err
	}

	for _, user := range users {
		// Optional: prevent CSV injection in Excel
		username := sanitizeCSV(user.Username)

		err := writer.Write([]string{
			username,
			user.UserDate,
			user.UpdateType,
			fmt.Sprintf("%d", user.DaysDelta),
		})
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	return &buf, nil
}

// Prevent Excel formula injection (optional but smart)
func sanitizeCSV(value string) string {
	if len(value) > 0 {
		switch value[0] {
		case '=', '+', '-', '@':
			return "'" + value
		}
	}
	return value
}

func (cron *ScheduledCronJob) findRoleMention(guildID, roleName string) string {
	roles, err := cron.session.GuildRoles(guildID)
	if err != nil {
		cron.logger.Warn("Failed to get guild roles", "error", err)
		return ""
	}

	for _, role := range roles {
		if role.Name == roleName {
			return fmt.Sprintf("<@&%s>", role.ID)
		}
	}

	cron.logger.Warn("Role not found", "roleName", roleName)
	return ""
}
