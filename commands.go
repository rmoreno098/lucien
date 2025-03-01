package main

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	vh  *VoiceHandler
	aqh *AudioQueueHandler
	res string

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
			Name:        "disconnect",
			Description: "Disconnect the bot from the voice channel",
		},
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
	}

	// Map commands to their respective handler functions
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"getusers": GetUsersHandler,
		// "search":   SearchHandler,
		"play":       PlayHandler,
		"disconnect": DisconnectHandler,
	}
)

func init() {
	vh = NewVoiceHandler()
	aqh = NewAudioQueueHandler()
}

func RegisterCommands(s *discordgo.Session) {
	for _, command := range Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, GUILD_ID, command)
		if err != nil {
			log.Fatalf("Error creating command: %v", err)
		}
	}
}

func PlayHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	GenerateResponse(s, i, discordgo.InteractionResponseDeferredChannelMessageWithSource, "")

	query := i.ApplicationCommandData().Options[0].StringValue()
	if query == "" {
		GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Please provide a valid song name or YouTube URL.")
		return
	}

	videoURL, err := resolveQuery(query)
	if err != nil {
		GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Could not find a valid YouTube video.")
		return
	}

	voiceConnection, err := vh.SetConnection(s)
	if err != nil {
		log.Printf("Error connecting to voice channel: %v", err)
		GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Error connecting to voice channel.")
		return
	}

	q := aqh.GetQueue()
	if vh.connections[GUILD_ID].IsPlaying || (len(q) > 0 && vh.connections[GUILD_ID].IsPlaying) {
		aqh.AddToQueue(videoURL)
		res = "Added " + videoURL + " to the queue."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &res,
		})
		return
	}

	res = "Now playing: " + videoURL
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &res,
	})

	if err := PlaySong(videoURL, voiceConnection, aqh, vh); err != nil {
		log.Printf("Playback error: %v", err)
		GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Error playing the song.")
		return
	}
}

func resolveQuery(query string) (string, error) {
	if strings.HasPrefix(query, "http://www.youtube.com") || strings.HasPrefix(query, "https://www.youtube.com") {
		youtubeRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/(watch\?v=|embed/|v/)?([a-zA-Z0-9_-]+)`)
		if youtubeRegex.MatchString(query) {
			return query, nil
		}
		return "", errors.New("invalid YouTube URL")
	}

	results, err := SearchYouTube(query)
	if err != nil {
		return "", err
	}

	return results[0], nil
}

func GetUsersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get list of users in the server
	guildID := i.GuildID
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		log.Printf("Error fetching members: %v", err)
		GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Error fetching members.")
		return
	}

	// Format member list
	var memberList []string
	for _, member := range members {
		memberList = append(memberList, member.User.Username)
	}

	// Respond with member list
	response := "Members:\n" + strings.Join(memberList, "\n")
	GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, response)
}

func DisconnectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	voiceConnection := vh.GetConnection(GUILD_ID)
	if voiceConnection != nil {
		aqh.ClearQueue()
		vh.Disconnect(s, GUILD_ID)
		GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Disconnected from voice channel.")
	}
	GenerateResponse(s, i, discordgo.InteractionResponseChannelMessageWithSource, "Not connected to any voice channel.")
}
