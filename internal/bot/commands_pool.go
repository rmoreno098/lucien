package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func NewCommandsPool() *WorkerPool {
	return &WorkerPool{
		tasks:   make(chan Task, 100),
		workers: 2,
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
		log.Printf("Removed guild from connections map: %v", guildID)
		return nil
	}
	log.Printf("Could not find voice connection for guild: %v", guildID)
	return nil
}
