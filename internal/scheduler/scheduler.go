package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduledCronJob struct {
	conn    *pgxpool.Pool
	logger  *slog.Logger
	session *discordgo.Session
}

type InactiveUser struct {
	Username   string
	UpdateType string
	UserDate   string
	DaysDelta  int
}

func NewScheduledCronJob(conn *pgxpool.Pool, logger *slog.Logger, session *discordgo.Session) *ScheduledCronJob {
	return &ScheduledCronJob{
		conn:    conn,
		logger:  logger,
		session: session,
	}
}

func (cron *ScheduledCronJob) GetUserActivityWithTimeRanges() {
	cron.logger.Info("Initializing scheduled cron job for fetching user activity time deltas in Discord")

	ctx := context.Background()

	// Get all user events from the database
	events, err := database.GetAllUserEvents(cron.conn, ctx, cron.logger)
	if err != nil {
		cron.logger.Error("Failed to fetch user events", "error", err)
		return
	}

	// Get today's date in "2006-01-02" format
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	todayDate, err := time.Parse("2006-01-02", todayStr)
	if err != nil {
		cron.logger.Error("Failed to parse today's date", "error", err)
		return
	}

	// List to store users that exceed 90 days
	var inactiveUsers []InactiveUser

	// Loop through each record and calculate the delta
	for _, event := range events {
		// Parse the user's date string
		userDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			cron.logger.Warn("Failed to parse user date", "username", event.Username, "date", event.Date, "error", err)
			continue
		}

		// Calculate the delta (difference) between today and the user's date
		delta := todayDate.Sub(userDate)
		daysDelta := int(delta.Hours() / 24)

		cron.logger.Info("User activity time delta calculated",
			"username", event.Username,
			"update_type", event.UpdateType,
			"user_date", event.Date,
			"today", todayStr,
			"days_delta", daysDelta,
		)

		// Add to inactive users list if delta exceeds 90 days
		if daysDelta > 90 {
			inactiveUsers = append(inactiveUsers, InactiveUser{
				Username:   event.Username,
				UpdateType: event.UpdateType,
				UserDate:   event.Date,
				DaysDelta:  daysDelta,
			})
		}
	}

	channelID := "1083587191336341654"
	// Send Discord message with inactive users if any found
	if len(inactiveUsers) > 0 && channelID != "" {
		// Get channel to find guild ID
		channel, err := cron.session.Channel(channelID)
		if err != nil {
			cron.logger.Error("Failed to get channel", "error", err)
			return
		}

		// Find Council role by name
		roleMention := cron.findRoleMention(channel.GuildID, "Council")

		message := cron.formatInactiveUsersMessage(inactiveUsers, todayStr)

		// Prepend role mention if found
		content := message
		if roleMention != "" {
			content = roleMention + "\n\n" + message
		}

		_, err = cron.session.ChannelMessageSend(channelID, content)
		if err != nil {
			cron.logger.Error("Failed to send Discord message", "error", err)
		} else {
			cron.logger.Info("Sent inactive users list to Discord", "count", len(inactiveUsers))
		}
	} else if len(inactiveUsers) > 0 && channelID == "" {
		cron.logger.Warn("Inactive users found but no channel ID configured", "count", len(inactiveUsers))
	}

	cron.logger.Info("Completed scheduled cron job for user activity time deltas", "total_records", len(events), "inactive_users", len(inactiveUsers))
}

func (cron *ScheduledCronJob) formatInactiveUsersMessage(users []InactiveUser, today string) string {
	var builder strings.Builder

	builder.WriteString("## Users Exceeding 90 Days Inactivity\n\n")
	builder.WriteString(fmt.Sprintf("**Report Date:** %s\n", today))
	builder.WriteString(fmt.Sprintf("**Total Users:** %d\n\n", len(users)))
	builder.WriteString(strings.Repeat("─", 40) + "\n\n")

	for i, user := range users {
		builder.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, user.Username))
		builder.WriteString(fmt.Sprintf("   Last Activity: %s\n", user.UserDate))
		builder.WriteString(fmt.Sprintf("   Update Type: %s\n", user.UpdateType))
		builder.WriteString(fmt.Sprintf("   Days Inactive: **%d**\n", user.DaysDelta))

		// Add separator between users (except for the last one)
		if i < len(users)-1 {
			builder.WriteString("\n" + strings.Repeat("─", 40) + "\n\n")
		}
	}

	return builder.String()
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
