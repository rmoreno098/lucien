package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	*WorkerPool
}

func NewCommandsPool() *Commands {
	return &Commands{
		WorkerPool: &WorkerPool{
			Tasks: make(chan Task, 100),
			handler: func(i Task) error {
				if err := i.Execute(); err != nil {
					return err
				}
				return nil
			},
			workers: 2,
		},
	}
}

type StopQueue struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Voice       *VoiceHandler
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
