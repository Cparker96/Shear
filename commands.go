package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/bwmarrin/discordgo"
// )

// func executeShowActivityMessage(s *discordgo.Session, m *discordgo.Message) {
// 	ctx := context.Background()
// 	pool := conn

// 	query := "SELECT * FROM public.activity"
// 	rows, err := pool.Query(ctx, query)
// 	if err != nil {
// 		logger.Error("Failed to query for recent activity", "Error", err)
// 		return
// 	}
// 	defer rows.Close()

// 	var responseText string
// 	var count int

// 	responseText += "**Activity Log:**\n"
// 	responseText += "```\n"
// 	responseText += fmt.Sprintf("%-5s | %-20s | %s\n", "Username", "Update_Type", "Date")
// 	responseText += "--------------------------------------------------------\n"

// 	for rows.Next() {
// 		var user string
// 		var updateType string
// 		var date string

// 		if err := rows.Scan(&user, &updateType, &date); err != nil {
// 			logger.Error("Error scanning row", "error", err)
// 			continue
// 		}

// 		parsedTime, err := time.Parse("2006-01-02", date)
// 		displayDate := date

// 		if err == nil {
// 			displayDate = parsedTime.Format("Jan 02, 2006")
// 		} else {
// 			logger.Warn("Could not parse date string", "Date", date, "Error", err)
// 		}

// 		responseText += fmt.Sprintf("%-10s | %-20s | %s\n", user, updateType, displayDate)
// 	}
// 	responseText += "```"

// 	if count == 0 {
// 		responseText = "Query successful, but no activities found in the table"
// 	}

// 	// Check for any error that occurred during row iteration
// 	if rows.Err() != nil {
// 		log.Printf("Error after iterating rows: %v", rows.Err())
// 		responseText = "Warning: Data retrieval error detected after some rows were read."
// 	}
// }

// func strPtr(s string) *string {
// 	return &s
// }
