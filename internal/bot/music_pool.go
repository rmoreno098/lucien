package bot

import (
	"fmt"
	"log"
	"lucien/internal/services"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type Music struct {
	*WorkerPool
	MusicMap []string
}

func NewMusicPool() *Music {
	return &Music{
		WorkerPool: &WorkerPool{
			Tasks: make(chan Task, 100),
			handler: func(i Task) error {
				if err := i.Execute(); err != nil {
					return err
				}
				return nil
			},
			workers: 1,
		},
		MusicMap: make([]string, 0),
	}
}

type PlaySong struct {
	Url         string
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Voice       *VoiceHandler
}

func (h *PlaySong) Execute() error {
	conn, err := h.Voice.SetConnection(h.Session, h.Interaction)
	if err != nil {
		h.Session.ChannelMessageSend(h.Interaction.ChannelID, err.Error())
		return err
	}
	defer h.cleanup(conn)

	audio, err := services.GetAudioURL(h.Url)
	if err != nil {
		log.Printf("An error occurred trying to get audio url: %v", err)
		h.Session.ChannelMessageSend(h.Interaction.ChannelID, err.Error())
		return err
	}

	h.Session.ChannelMessageSend(h.Interaction.ChannelID, fmt.Sprintf("Now playing: %s", h.Url))
	log.Printf("Now playing: %v", h.Url)

	done := make(chan bool)
	dgvoice.PlayAudioFile(conn, audio, done)
	<-done

	return nil
}

func (h *PlaySong) cleanup(connection *discordgo.VoiceConnection) {
	h.Voice.mu.Lock()
	defer h.Voice.mu.Unlock()

	connection.Close()

	v := h.Voice.connections[h.Interaction.GuildID]
	v.IsConnected = false
	v.IsPlaying = false

	delete(h.Voice.connections, h.Interaction.GuildID)
}
