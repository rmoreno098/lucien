package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

// GetAudioURL fetches the audio stream URL using yt-dlp.
func GetAudioURL(videoURL string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--get-url", videoURL)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func PlaySong(videoURL string, vc *discordgo.VoiceConnection, aqh *AudioQueueHandler, vh *VoiceHandler) error {
	if vc == nil || vc.ChannelID == "" {
		log.Println("Voice connection is nil or not properly initialized.")
		return fmt.Errorf("voice connection is nil or not properly initialized")
	}

	var audio = videoURL

	if len(aqh.GetQueue()) > 0 {
		audio = aqh.RemoveFromQueue().URL
	}

	streamURL, err := GetAudioURL(audio)
	if err != nil {
		log.Printf("Error fetching audio stream URL: %v\n", err)
		return err
	}

	vh.connections[GUILD_ID].IsPlaying = true
	dgvoice.PlayAudioFile(vc, streamURL, make(chan bool))
	vh.connections[GUILD_ID].IsPlaying = false

	return nil
}
