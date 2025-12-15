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
	conn          *pgxpool.Pool
	ctx           context.Context
	logger        *slog.Logger
	CommandPrefix = "!shear"
)

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	token, postgresPW, postgresUser, postgresDBName, err := GetEnvValues()
	conn, ctx, _ = PostgresConn(postgresUser, postgresPW, postgresDBName)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Error creating Discord session:", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsGuildMessageReactions

	dg.AddHandler(ready)
	dg.AddHandler(VoiceState)
	dg.AddHandler(MessageCreate)
	dg.AddHandler(MessageReaction)
	dg.AddHandler(MessageCreateFromCommand)

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
