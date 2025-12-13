package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// VoiceStateUpdate runs automatically whenever someone's voice state changes
// this function runs indefinitely as long as the bot is connected
func VoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Skip initial state
	if vs.BeforeUpdate == nil {
		return
	}

	// Get user info
	member, err := s.GuildMember(vs.GuildID, vs.UserID)
	if err != nil {
		log.Printf("Error getting member: %v", err)
		return
	}

	// Ignore bots
	if member.User.Bot {
		return
	}

	// Get guild info
	guild, err := s.Guild(vs.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return
	}

	oldChannelID := vs.BeforeUpdate.ChannelID
	newChannelID := vs.ChannelID

	// USER JOINED VOICE CHANNEL
	if oldChannelID == "" && newChannelID != "" {
		channelName := getChannelName(guild, newChannelID)
		fmt.Printf("[VOICE JOIN] %s entered %s at %s\n",
			member.User.Username,
			channelName,
			time.Now().Format("15:04:05"))

		// ADD YOUR LOGIC HERE
		handleUserJoined(vs.UserID, member.User.Username, channelName)
	}

	// USER LEFT VOICE CHANNEL
	if oldChannelID != "" && newChannelID == "" {
		channelName := getChannelName(guild, oldChannelID)
		fmt.Printf("[VOICE LEAVE] %s left %s at %s\n",
			member.User.Username,
			channelName,
			time.Now().Format("15:04:05"))

		// ADD YOUR LOGIC HERE
		handleUserLeft(vs.UserID, member.User.Username, channelName)
	}

	// USER SWITCHED CHANNELS
	if oldChannelID != "" && newChannelID != "" && oldChannelID != newChannelID {
		oldName := getChannelName(guild, oldChannelID)
		newName := getChannelName(guild, newChannelID)
		fmt.Printf("[VOICE SWITCH] %s: %s → %s at %s\n",
			member.User.Username,
			oldName,
			newName,
			time.Now().Format("15:04:05"))

		// ADD YOUR LOGIC HERE
		handleUserSwitched(vs.UserID, member.User.Username, oldName, newName)
	}
}

// Helper function to get channel name
func getChannelName(guild *discordgo.Guild, channelID string) string {
	for _, channel := range guild.Channels {
		if channel.ID == channelID {
			return channel.Name
		}
	}
	return "Unknown"
}

// YOUR CUSTOM LOGIC FUNCTIONS - Put your code here!

func handleUserJoined(userID, username, channelName string) {
	// Example: Save to database
	// db.Exec("INSERT INTO voice_joins ...")

	// Example: Send notification
	// sendNotification(username + " joined " + channelName)

	log.Printf("CUSTOM LOGIC: User %s joined %s", username, channelName)
}

func handleUserLeft(userID, username, channelName string) {
	// Example: Calculate session duration
	// duration := time.Since(joinTimes[userID])

	log.Printf("CUSTOM LOGIC: User %s left %s", username, channelName)
}

func handleUserSwitched(userID, username, oldChannel, newChannel string) {
	log.Printf("CUSTOM LOGIC: User %s switched from %s to %s", username, oldChannel, newChannel)
}
