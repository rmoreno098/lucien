package main

import (
	"github.com/bwmarrin/discordgo"
)

func GenerateResponse(s *discordgo.Session, i *discordgo.InteractionCreate, t discordgo.InteractionResponseType, m string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: t,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
}
