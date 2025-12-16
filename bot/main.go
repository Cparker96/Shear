package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/codyw/shear/internal/command"
	"github.com/codyw/shear/internal/config"
	"github.com/codyw/shear/internal/database"
	"github.com/codyw/shear/internal/event"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn          *pgxpool.Pool
	ctx           context.Context
	logger        *slog.Logger
	CommandPrefix = "!shear"
)

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	token, postgresPW, postgresUser, postgresDBName, err := config.GetEnvValues(logger)
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

	err = dg.Open()
	if err != nil {
		logger.Error("Error opening connection:", err)
	}
	defer dg.Close()

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
}

func CommandRouter(s *discordgo.Session, message *discordgo.MessageCreate, conn *pgxpool.Pool, logger *slog.Logger) {
	event.HandleShearCommand(s, message, conn, CommandPrefix, "get-activity", command.ExecuteGetUserActivity, logger)
	event.HandleShearCommand(s, message, conn, CommandPrefix, "remove-user", command.ExecuteRemoveUser, logger)
}
