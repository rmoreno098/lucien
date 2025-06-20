package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	DiscordSecret    string
	DiscordGuild     string
	DiscordChannelID string
}

func LoadConfig() Configuration {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return Configuration{
		DiscordSecret:    getEnv("DISCORD_TOKEN", ""),
		DiscordGuild:     getEnv("DISCORD_GUILD", ""),
		DiscordChannelID: getEnv("DISCORD_CHANNEL", ""),
	}
}

func getEnv(k, f string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}

	if f == "" {
		panic("Environment variable " + k + " not defined")
	}

	return f
}
