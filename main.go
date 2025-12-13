package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn   *pgxpool.Pool
	ctx    context.Context
	logger *slog.Logger
)

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	_, pw, _ := GetEnvValues(logger)
	conn, ctx, _ = PostgresConn(pw, logger)

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		logger.Error("BOT_TOKEN environment variable not set")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Error creating Discord session:", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildVoiceStates

	dg.AddHandler(ready)
	dg.AddHandler(VoiceStateUpdate)

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
