package bot

import (
	"errors"
	"log"
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

	h.Queue.AddToQueue(url, s, i, h.VoiceState)

	res := "Track has been added to the queue"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &res,
	})
}

func (h *DiscordHandler) DisconnectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.Queue.EndQueue(s, i, h.VoiceState)

	log.Printf("Disconnected from %s", i.GuildID)
	utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Disconnected from voice channel.")
}

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
