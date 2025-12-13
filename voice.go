package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Skip initial state
	if vs.BeforeUpdate == nil {
		return
	}

	isInVoice := vs.ChannelID != ""
	channelChanged := vs.BeforeUpdate.ChannelID != vs.ChannelID

	if !isInVoice && channelChanged {
		// Get user info
		member, err := s.GuildMember(vs.GuildID, vs.UserID)
		if err != nil || member.User.Bot {
			return
		}

		// Get guild info
		guild, err := s.Guild(vs.GuildID)
		if err != nil {
			logger.Error("Error getting guild", "Error", err)
			return
		}

		channelName := getChannelName(guild, vs.ChannelID)
		updatedTime := time.Now().Format("2006-01-02")
		logger.Info("User joined a voice channel", "User", member.User.Username, "Channel", channelName, "Time", updatedTime)

		go WriteToPostgres(conn, ctx, "voice", updatedTime, member.User.Username)
	}
}

func getChannelName(guild *discordgo.Guild, channelID string) string {
	for _, channel := range guild.Channels {
		if channel.ID == channelID {
			return channel.Name
		}
	}
	return "Unknown channel name"
}
