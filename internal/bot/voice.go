package bot

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type VoiceConnectionEntry struct {
	VoiceConnection *discordgo.VoiceConnection
	IsConnected     bool
	IsPlaying       bool
}

type VoiceHandler struct {
	mu          sync.Mutex
	connections map[string]*VoiceConnectionEntry // map of bot voice connections (per-guild)
}

func NewVoiceHandler() *VoiceHandler {
	return &VoiceHandler{
		connections: make(map[string]*VoiceConnectionEntry),
	}
}

func (h *VoiceHandler) SetConnection(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.VoiceConnection, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	guildID := i.GuildID
	c, exists := h.connections[guildID]
	if exists && c.IsConnected {
		log.Printf("Found existing connection for guildID %v", guildID)
		c.IsConnected = true
		c.IsPlaying = true
		return c.VoiceConnection, nil
	}

	channelID, err := getUserVoiceChannelID(s, guildID, i.Member.User.ID)
	if err != nil {
		log.Printf("Error getting user voice channel: %v", err)
		return nil, err
	}

	conn, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Printf("Error joining voice channel: %v", err)
		return nil, err
	}

	h.connections[guildID] = &VoiceConnectionEntry{
		VoiceConnection: conn,
		IsConnected:     true,
		IsPlaying:       true,
	}
	return conn, nil
}

func (h *VoiceHandler) GetConnection(i *discordgo.InteractionCreate) *VoiceConnectionEntry {
	return h.connections[i.GuildID]
}

func getUserVoiceChannelID(s *discordgo.Session, guildID, userID string) (string, error) {
	vs, err := s.State.VoiceState(guildID, userID)
	if err != nil || vs == nil {
		return "", fmt.Errorf("User not in voice channel")
	}
	return vs.ChannelID, nil
}
