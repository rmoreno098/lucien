package main

import (
	"log"
	config "lucien"
	"lucien/internal/bot"
	"lucien/pkg/discord"
)

func main() {
	c := config.LoadConfig()

	d := discord.NewSession(c.DiscordSecret, c.DiscordGuild, c.DiscordChannelID)
	if err := bot.RegisterDiscordSession(d); err != nil {
		log.Fatalf("Unable to launch bot: %v", err)
	}

	log.Printf("Running bot...")
	select {}
}
