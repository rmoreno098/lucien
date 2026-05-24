package bot

import (
	"fmt"
	"log"
	"lucien/internal/services"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func NewMusicPool() *WorkerPool {
	return &WorkerPool{
		tasks:   make(chan Task, 100),
		workers: 1,
	}
}

type PlaySongTask struct {
	Url         string
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Voice       *VoiceHandler
}

func (h *PlaySongTask) Execute() error {
	conn, err := h.Voice.SetConnection(h.Session, h.Interaction)
	if err != nil {
		h.Session.ChannelMessageSend(h.Interaction.ChannelID, err.Error())
		return err
	}

	log.Println("Connection set")
	defer h.cleanup(conn)

	audio, err := services.GetAudioURL(h.Url)
	if err != nil {
		log.Printf("An error occurred trying to get audio url: %v", err)
		h.Session.ChannelMessageSend(h.Interaction.ChannelID, err.Error())
		return err
	}

	log.Printf("Now playing: %v", h.Url)
	_, err = h.Session.ChannelMessageSend(h.Interaction.ChannelID, fmt.Sprintf("Now playing: %s", h.Url))
	if err != nil {
		return err
	}

	done := make(chan bool)
	dgvoice.PlayAudioFile(conn, audio, done)
	<-done

	return nil
}

func (h *PlaySongTask) cleanup(connection *discordgo.VoiceConnection) {
	h.Voice.mu.Lock()
	defer h.Voice.mu.Unlock()

	connection.Close()

	_, exists := h.Voice.connections[h.Interaction.GuildID]
	if exists {
		delete(h.Voice.connections, h.Interaction.GuildID)
	}
}
