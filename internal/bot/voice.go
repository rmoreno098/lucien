package bot

import (
	"fmt"
	"log"
	"lucien/internal/services"
	"sync"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type VoiceConnectionEntry struct {
	VoiceConnection *discordgo.VoiceConnection
	IsConnected     bool
	IsPlaying       bool
}

type VoiceHandler struct {
	mu          sync.Mutex
	connections map[string]*VoiceConnectionEntry
}

func NewVoiceHandler() *VoiceHandler {
	return &VoiceHandler{
		connections: make(map[string]*VoiceConnectionEntry),
	}
}

func (h *VoiceHandler) Play(request *TrackRequest) error {
	conn, err := h.SetConnection(request.Session, request.Interaction)
	if err != nil {
		return err
	}

	audio, err := services.GetAudioURL(request.Url)
	if err != nil {
		return err
	}

	done := make(chan bool)
	dgvoice.PlayAudioFile(conn, audio, done)
	<-done

	return nil
}

func (h *VoiceHandler) SetConnection(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.VoiceConnection, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	guildID := i.GuildID
	channelID, err := getUserVoiceChannelID(s, guildID, i.Member.User.ID)
	if err != nil {
		return nil, err
	}

	connEntry, exists := h.connections[guildID]
	if exists && connEntry.IsConnected {
		return connEntry.VoiceConnection, nil
	}

	conn, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	h.connections[guildID] = &VoiceConnectionEntry{
		VoiceConnection: conn,
		IsConnected:     true,
		IsPlaying:       true,
	}

	return conn, nil
}

func (h *VoiceHandler) getConnection(guildID string) *discordgo.VoiceConnection {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.connections[guildID]; ok {
		return conn.VoiceConnection
	}
	return nil
}

func getUserVoiceChannelID(s *discordgo.Session, guildID, userID string) (string, error) {
	vs, err := s.State.VoiceState(guildID, userID)
	if err != nil || vs == nil {
		return "", fmt.Errorf("user not in voice channel")
	}
	return vs.ChannelID, nil
}

func (h *VoiceHandler) Disconnect(s *discordgo.Session, guildID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, ok := h.connections[guildID]; ok && conn.IsConnected {
		if err := conn.VoiceConnection.Disconnect(); err != nil {
			log.Fatalf("Unable to disconnect voice connection: %v", err)
		}
		delete(h.connections, guildID)
	}
}
