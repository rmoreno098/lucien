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

	// Create a new Discord session
	var err error
	discord, err = discordgo.New("Bot " + API_KEY)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	discord.LogLevel = discordgo.LogWarning

	// Register the interaction handler for slash commands
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Received interaction: %s", i.ApplicationCommandData().Name)
		// Check if a handler exists for the invoked command
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i) // Call the handler function if found
		}
	})
}

func main() {
	// Open a WebSocket connection to Discord
	err := discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	RegisterCommands(discord)

	defer discord.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Wait for interrupt signal
	log.Println("Bot is running. Press Ctrl+C to exit.")
	<-stop

	log.Println("Shutting down...")
}
