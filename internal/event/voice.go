package event

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func VoiceState(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
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

		updatedTime := time.Now().Format("2006-01-02")
		logger.Info("User joined a voice channel", "User", member.User.Username, "Time", updatedTime)

		go WriteToPostgres(conn, ctx, "voice", updatedTime, member.User.Username)
	}
}
