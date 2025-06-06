package main

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type VoiceConnectionState struct {
	VoiceConnection *discordgo.VoiceConnection
	IsConnected     bool
	IsPlaying       bool
	GuildID         string
	ChannelID       string
}

type VoiceHandler struct {
	connections map[string]*VoiceConnectionState
	mu          sync.Mutex
}

func NewVoiceHandler() *VoiceHandler {
	return &VoiceHandler{
		connections: make(map[string]*VoiceConnectionState),
	}
}

func (vh *VoiceHandler) GetConnection(guildID string) *discordgo.VoiceConnection {
	vh.mu.Lock()
	defer vh.mu.Unlock()

	if conn, ok := vh.connections[guildID]; ok {
		return conn.VoiceConnection
	}
	return nil
}

func (vh *VoiceHandler) SetConnection(s *discordgo.Session) (*discordgo.VoiceConnection, error) {
	vh.mu.Lock()
	defer vh.mu.Unlock()

	// Check if the bot is already connected to a voice channel
	if conn, ok := vh.connections[GUILD_ID]; ok && conn.IsConnected {
		log.Println("Bot is already connected to a voice channel.")
		return conn.VoiceConnection, nil
	}

	vc, err := s.ChannelVoiceJoin(GUILD_ID, CHANNEL_ID, false, true)
	if err != nil {
		log.Printf("Error connecting to voice channel: %v\n", err)
		return nil, err
	}

	vh.connections[GUILD_ID] = &VoiceConnectionState{
		VoiceConnection: vc,
		IsConnected:     true,
		GuildID:         GUILD_ID,
		ChannelID:       CHANNEL_ID,
	}

	return vc, nil
}

func (vh *VoiceHandler) Disconnect(s *discordgo.Session, guildID string) {
	vh.mu.Lock()
	defer vh.mu.Unlock()

	if conn, ok := vh.connections[guildID]; ok && conn.IsConnected {
		conn.VoiceConnection.Disconnect()
		conn.IsConnected = false
		delete(vh.connections, guildID)
		log.Println("Disconnected from voice channel.")
	}
}
