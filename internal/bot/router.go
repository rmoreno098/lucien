package bot

import (
	"log"
	"lucien/pkg/discord"

	"github.com/bwmarrin/discordgo"
)

type DiscordRegister struct {
	Session *discordgo.Session
}

func RegisterDiscordSession(service *discord.DiscordService) error {
	s := registerServices()

	registerCommands(service)
	registerHandlers(service.Session, s)
	if err := service.Session.Open(); err != nil {
		return err
	}

	return nil
}

func registerServices() *DiscordHandler {
	voice := NewVoiceHandler() // Manages a mapping of voice connections by guild
	queue := NewAudioQueue()   // Manages a worker pool of threads to manage user requests

	queue.StartWorkers(2, func(request *TrackRequest) error {
		return voice.Play(request)
	})

	return &DiscordHandler{
		Queue:      queue,
		VoiceState: voice,
	}
}

func registerHandlers(s *discordgo.Session, h *DiscordHandler) {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			h.PlayHandler(s, i)
		},
		"disconnect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			h.DisconnectHandler(s, i)
		},
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name
		if handler, ok := commandHandlers[name]; ok {
			handler(s, i)
		}
	})
}

func registerCommands(s *discord.DiscordService) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Command to play a song from YouTube",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "song",
					Description: "Provide Search Query or YouTube URL",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "disconnect",
			Description: "Disconnect the bot from the voice channel",
		},
		{
			Name:        "getusers",
			Description: "Get a list of users in the server",
		},
		{
			Name:        "search",
			Description: "Command to search for a song on YouTube",
		},
	}
	for _, command := range commands {
		_, err := s.Session.ApplicationCommandCreate(s.AppID, s.Guild, command)
		if err != nil {
			log.Fatalf("Error creating command: %v", err)
		}
	}
}
