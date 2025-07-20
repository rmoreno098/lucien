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
	user_query := i.ApplicationCommandData().Options[0].StringValue()

	url, err := resolveQuery(user_query)
	if err != nil {
		utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Could not find a valid YouTube video.")
		return
	}

	h.Queue.MusicPool.AddTask(url, s, i, h.VoiceState)

	res := "Track has been added to the queue"
	utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, res)
}

func (h *DiscordHandler) DisconnectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.Queue.CommandsPool.EndQueue(s, i, h.VoiceState)

	log.Printf("Disconnected from %s", i.GuildID)
	utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Disconnected from voice channel.")
}

func (h *DiscordHandler) QueueStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	list := fmt.Sprintf("h.Queue.MusicPool.MusicMap: %v\n", h.Queue.MusicPool.MusicMap)
	utils.GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, list)
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
