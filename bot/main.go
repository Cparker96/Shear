package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/command"
	"github.com/codyw/shear/internal/config"
	"github.com/codyw/shear/internal/database"
	"github.com/codyw/shear/internal/event"
	"github.com/codyw/shear/internal/scheduler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

var (
	conn          *pgxpool.Pool
	ctx           context.Context
	logger        *slog.Logger
	CommandPrefix = "!shear"
)

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	token, postgresPW, postgresUser, postgresDBName, channelID, roleName, err := config.GetEnvValues(logger)
	if err != nil {
		logger.Error("Failed to load environment variables", "error", err)
		return
	}
	conn, ctx, _ = database.PostgresConn(postgresUser, postgresPW, postgresDBName, logger)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Error creating Discord session:", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentGuildMembers

	dg.AddHandler(ready)

	dg.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
		event.VoiceState(s, vs, conn, logger)
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		event.MessageCreate(s, m, conn, logger)
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		event.MessageReaction(s, r, conn, logger)
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		CommandRouter(s, m, conn, logger)
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
		event.DetectUserJoin(s, r, conn, logger)
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.GuildMemberRemove) {
		event.DetectUserLeave(s, r, conn, logger)
	})

	err = dg.Open()
	if err != nil {
		logger.Error("Error opening connection:", err)
	}
	defer dg.Close()

	// Initialize and start the cron job scheduler
	cronJob := scheduler.NewScheduledCronJob(conn, logger, dg, channelID, roleName)

	// Set up cron scheduler to run nightly at 2 AM
	c := cron.New()
	_, err = c.AddFunc("0 2 * * *", cronJob.GetUserActivityWithTimeRanges)
	if err != nil {
		logger.Error("Failed to schedule cron job", "error", err)
	} else {
		logger.Info("Scheduled cron job to run nightly at 2 AM")
		c.Start()
		defer c.Stop()
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// this keeps the program running forever (until CTRL-C)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc // blocks here forever, keeping bot alive

	fmt.Println("\nShutting down...")
}

// Called when bot connects
func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot connected successfully")
	fmt.Printf("Bot User: %s\n", s.State.User.Username)
	fmt.Println("\nMonitoring voice channels in:")
	for _, guild := range event.Guilds {
		g, _ := s.Guild(guild.ID)
		fmt.Printf("  → %s\n", g.Name)
	}
	fmt.Println()

	// Seed database with all guild members if table is empty
	seedDatabaseIfEmpty(s, conn, ctx, logger)
}

func seedDatabaseIfEmpty(s *discordgo.Session, conn *pgxpool.Pool, ctx context.Context, logger *slog.Logger) {
	// Check if database has any records
	hasRecords, err := database.HasRecords(conn, ctx, logger)
	if err != nil {
		logger.Error("Failed to check if database has records", "error", err)
		return
	}

	if hasRecords {
		logger.Info("Database already contains records, skipping seed")
		return
	}

	logger.Info("Database is empty, starting seed process")

	// Get all guilds the bot is in
	guilds := s.State.Guilds
	if len(guilds) == 0 {
		logger.Warn("Bot is not in any guilds, cannot seed database")
		return
	}

	// Collect all unique usernames from all guilds
	usernamesMap := make(map[string]bool)

	for _, guild := range guilds {
		logger.Info("Fetching members from guild", "guild_id", guild.ID, "guild_name", guild.Name)

		// Fetch all members (Discord paginates at 1000, so we need to handle pagination)
		after := ""
		for {
			members, err := s.GuildMembers(guild.ID, after, 1000)
			if err != nil {
				logger.Error("Failed to fetch guild members", "guild_id", guild.ID, "error", err)
				break
			}

			if len(members) == 0 {
				break
			}

			for _, member := range members {
				// Skip bots
				if member.User.Bot {
					continue
				}
				// Use username (not nickname) as the identifier
				usernamesMap[member.User.Username] = true
			}

			// Check if we got fewer than 1000 members (last page)
			if len(members) < 1000 {
				break
			}

			// Set after to the last member's ID for pagination
			after = members[len(members)-1].User.ID
		}
	}

	// Convert map to slice
	usernames := make([]string, 0, len(usernamesMap))
	for username := range usernamesMap {
		usernames = append(usernames, username)
	}

	logger.Info("Collected usernames for seeding", "count", len(usernames))

	// Get today's date in "2006-01-02" format
	today := time.Now().Format("2006-01-02")

	// Seed the database
	err = database.SeedDatabase(conn, ctx, usernames, today, logger)
	if err != nil {
		logger.Error("Failed to seed database", "error", err)
		return
	}

	logger.Info("Database seeding completed successfully")
}

func CommandRouter(s *discordgo.Session, message *discordgo.MessageCreate, conn *pgxpool.Pool, logger *slog.Logger) {
	event.HandleShearCommand(s, message, conn, CommandPrefix, "get-activity", command.ExecuteGetUserActivity, logger)
	event.HandleShearCommand(s, message, conn, CommandPrefix, "remove-user", command.ExecuteRemoveUser, logger)
}
