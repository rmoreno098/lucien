package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	discord    *discordgo.Session
	API_KEY    string
	GUILD_ID   string
	CHANNEL_ID string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	API_KEY = os.Getenv("DISCORD_TOKEN")
	GUILD_ID = os.Getenv("DISCORD_GUILD")
	CHANNEL_ID = os.Getenv("DISCORD_CHANNEL")

	var err error
	discord, err = discordgo.New("Bot " + API_KEY)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	discord.LogLevel = discordgo.LogWarning

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Received interaction: %s", i.ApplicationCommandData().Name)
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	err := discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	RegisterCommands(discord)

	defer discord.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Println("Bot is running. Press Ctrl+C to exit.")
	<-stop

	log.Println("Shutting down...")
}
