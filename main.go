package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	_, pw, _ := GetEnvValues(logger)
	conn, ctx, _ := PostgresConn(pw, logger)

	a, _ := FindData(conn, ctx, logger)
	fmt.Println(a)
	// 	token := os.Getenv("BOT_TOKEN")
	// 	if token == "" {
	// 		logger.Error("BOT_TOKEN environment variable not set")
	// 	}

	// 	dg, err := discordgo.New("Bot " + token)
	// 	if err != nil {
	// 		logger.Error("Error creating Discord session:", err)
	// 	}

	// 	dg.Identify.Intents = discordgo.IntentsGuilds |
	// 		discordgo.IntentsGuildMessages |
	// 		discordgo.IntentsMessageContent |
	// 		discordgo.IntentsGuildVoiceStates

	// 	// Register handlers
	// 	dg.AddHandler(ready)
	// 	// dg.AddHandler(messageCreate)
	// 	// dg.AddHandler(voiceStateUpdate)

	// 	err = dg.Open()
	// 	if err != nil {
	// 		logger.Error("Error opening connection:", err)
	// 	}
	// 	defer dg.Close()

	// 	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// 	// This keeps the program running forever (until CTRL-C)
	// 	sc := make(chan os.Signal, 1)
	// 	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// 	<-sc // Blocks here forever, keeping bot alive

	// 	fmt.Println("\nShutting down...")
	// }

	// // Called when bot connects
	// func ready(s *discordgo.Session, event *discordgo.Ready) {
	// 	fmt.Println("Bot connected successfully")
	// 	fmt.Printf("Bot User: %s\n", s.State.User.Username)
	// 	fmt.Println("\nMonitoring voice channels in:")
	// 	for _, guild := range event.Guilds {
	// 		g, _ := s.Guild(guild.ID)
	// 		fmt.Printf("  → %s\n", g.Name)
	// 	}
	// 	fmt.Println()
}
