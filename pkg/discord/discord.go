package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordService struct {
	Session   *discordgo.Session
	Guild     string
	ChannelID string
}

func NewSession(secret, guild, channel string) *DiscordService {
	session, err := discordgo.New("Bot " + secret)
	if err != nil {
		log.Fatalf("An error occurred initalizing Discord service: %v", err)
	}

	return &DiscordService{
		Session:   session,
		Guild:     guild,
		ChannelID: channel,
	}
}
