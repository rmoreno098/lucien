package bot

import (
	"errors"
	"fmt"
	"lucien/pkg/utils"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordHandler struct {
	Queue      *AudioQueue
	VoiceState *VoiceHandler
}

func (h *DiscordHandler) PlayHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	utils.GenerateResponse(s, i, discordgo.InteractionResponseDeferredChannelMessageWithSource, "")
	q := i.ApplicationCommandData().Options[0].StringValue()

	url, err := resolveQuery(q)
	if err != nil {
		utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Could not find a valid YouTube video.")
		return
	}

	h.Queue.AddToQueue(url, s, i)

	res := fmt.Sprintf("%s has been added to the queue", url)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &res,
	})
}

// func GetUsersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	// // Get list of users in the server
// 	// guildID := i.GuildID
// 	// members, err := s.GuildMembers(guildID, "", 1000)
// 	// if err != nil {
// 	// 	log.Printf("Error fetching members: %v", err)
// 	// 	utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Error fetching members.")
// 	// 	return
// 	// }

// 	// // Format member list
// 	// var memberList []string
// 	// for _, member := range members {
// 	// 	memberList = append(memberList, member.User.Username)
// 	// }

// 	// // Respond with member list
// 	// response := "Members:\n" + strings.Join(memberList, "\n")
// 	// utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, response)
// }

// func DisconnectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	voiceConnection := vh.GetConnection(GUILD_ID)
// 	if voiceConnection != nil {
// 		aqh.ClearQueue()
// 		vh.Disconnect(s, GUILD_ID)
// 		utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Disconnected from voice channel.")
// 	}
// 	utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Not connected to any voice channel.")
// }

func resolveQuery(query string) (string, error) {
	if strings.HasPrefix(query, "http://www.youtube.com") || strings.HasPrefix(query, "https://www.youtube.com") {
		youtubeRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/(watch\?v=|embed/|v/)?([a-zA-Z0-9_-]+)`)
		if youtubeRegex.MatchString(query) {
			return query, nil
		}
		return "", errors.New("invalid YouTube URL")
	}
	return "", errors.New("fix me")
}
