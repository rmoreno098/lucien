package main

import (
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

var (
	// YouTube client
	YTClient youtube.Client

	// Provide params for commands
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "getusers",
			Description: "Get a list of users in the server",
		},
		{
			Name:        "search",
			Description: "Command to search for a song on YouTube",
		},
		{
			Name:        "play",
			Description: "Command to play a song from YouTube",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "song",
					Description: "Query or YouTube URL",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	// Map commands to their respective handler functions
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"getusers": GetUsersHandler,
		"search":   SearchHandler,
		"play":     PlayHandler,
	}
)

func init() {
	YTClient = youtube.Client{}
}

func RegisterCommands(s *discordgo.Session) {
	for _, command := range Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, GUILD_ID, command)
		if err != nil {
			log.Fatalf("Error creating command: %v", err)
		}
	}

	log.Println("Commands registered.")
}

func SearchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	video, err := YTClient.GetVideo("https://www.youtube.com/watch?v=TS0C_pt08Go&t=1584s")
	if err != nil {
		log.Printf("An error occurred trying to get your video: %s", err)
	}
	log.Println("Video:", video)
}

func PlayHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var query string
	for _, options := range i.ApplicationCommandData().Options {
		if options.Name == "song" {
			query = options.StringValue()
		}
	}

	if query == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a song name or URL to play.",
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Now playing: " + query,
		},
	})

	voiceConnection, err := s.ChannelVoiceJoin(GUILD_ID, CHANNEL_ID, false, false)
	if err != nil {
		log.Printf("An error occurred trying to join the channel: %v", err)
	}

	voiceConnection.LogLevel = discordgo.LogWarning

	for voiceConnection.Ready == false {
		runtime.Gosched()
	}

	time.Sleep(time.Millisecond * 250)

	// Check if the query is a url or a query
	if strings.HasPrefix(query, "http://www.youtube.com") || strings.HasPrefix(query, "https://www.youtube.com") {
		PlaySong(query, voiceConnection)
	}

	voiceConnection.Close()
}

func GetUsersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("Here are the users")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "test",
		},
	})
}
