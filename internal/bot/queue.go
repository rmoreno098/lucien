package bot

import (
	"fmt"
	"log"
	"lucien/internal/services"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type Job interface {
	Execute() error
}

type WorkerPool struct {
	Jobs chan Job
	Quit chan bool
}

type AudioQueue struct {
	Pool *WorkerPool
}

func NewAudioQueue() *AudioQueue {
	pool := &WorkerPool{
		Jobs: make(chan Job, 100),
	}
	pool.startWorkers(2)

	return &AudioQueue{
		Pool: pool,
	}
}

func (wp *WorkerPool) startWorkers(numWorkers int) {
	for i := range numWorkers {
		go func(id int) {
			for {
				select {
				case job := <-wp.Jobs:
					if err := job.Execute(); err != nil {
						log.Printf("An error occurred trying to perform task: %v", err)
					}
				}
			}
		}(i)
	}
}

type PlaySong struct {
	Url         string
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Voice       *VoiceHandler
}

func (q *AudioQueue) AddToQueue(u string, s *discordgo.Session, i *discordgo.InteractionCreate, v *VoiceHandler) {
	q.Pool.Jobs <- &PlaySong{
		Url:         u,
		Session:     s,
		Interaction: i,
		Voice:       v,
	}
}

func (h *PlaySong) Execute() error {
	conn, err := h.Voice.SetConnection(h.Session, h.Interaction)
	if err != nil {
		h.Session.ChannelMessageSend(h.Interaction.ChannelID, err.Error())
		return err
	}

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

type StopQueue struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Voice       *VoiceHandler
}

func (q *AudioQueue) EndQueue(s *discordgo.Session, i *discordgo.InteractionCreate, v *VoiceHandler) {
	q.Pool.Jobs <- &StopQueue{
		Session:     s,
		Interaction: i,
		Voice:       v,
	}
}

func (h *StopQueue) Execute() error {
	h.Voice.mu.Lock()
	defer h.Voice.mu.Unlock()

	guildID := h.Interaction.GuildID
	if conn, ok := h.Voice.connections[guildID]; ok {
		err := conn.VoiceConnection.Disconnect()
		if err != nil {
			log.Printf("An error occurred trying to disconnect connection: %v", err)
			return err
		}
		delete(h.Voice.connections, guildID)
		log.Printf("Disconnected from voice channel for guild: %v", guildID)
		return nil
	}
	log.Printf("Could not find voice connection for guild: %v", guildID)
	return nil
}

func (q *AudioQueue) Shutdown() {
	close(q.Pool.Jobs)
}
