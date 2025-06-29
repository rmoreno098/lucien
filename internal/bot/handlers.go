package bot

import (
	"errors"
	"fmt"
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
	if h.Queue == nil {
		utils.GenerateResponse(s, i, discordgo.InteractionResponseDeferredChannelMessageWithSource, "There is no queue!!!")
	}

	utils.GenerateResponse(s, i, discordgo.InteractionResponseDeferredChannelMessageWithSource, "")
	q := i.ApplicationCommandData().Options[0].StringValue()

	url, err := resolveQuery(q)
	if err != nil {
		utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Could not find a valid YouTube video.")
		return
	}

	h.Queue.AddToQueue(url, s, i)

	res := fmt.Sprintf("Track has been added to the queue")
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &res,
	})
}

func (h *DiscordHandler) DisconnectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID

	if err := h.VoiceState.Disconnect(guildID); err != nil {
		log.Printf("There was an error disconnecting: %v", err)
		utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, err.Error())
	}
	h.Queue.Quit <- true

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
