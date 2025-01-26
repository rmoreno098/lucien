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

// VideoInfo represents YouTube video information fetched by yt-dlp.
type VideoInfo struct {
	URL string `json:"url"`
}

// GetAudioURL fetches the audio stream URL using yt-dlp.
func GetAudioURL(videoURL string) (string, error) {
	// yt-dlp command to get the best audio stream URL in JSON format
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "--get-url", videoURL)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Return the first valid URL (trim any extra spaces/newlines)
	return strings.TrimSpace(out.String()), nil
}

func PlaySong(videoURL string, vc *discordgo.VoiceConnection) error {
	if vc == nil || vc.ChannelID == "" {
		log.Println("Voice connection is nil or not properly initialized.")
		return fmt.Errorf("voice connection is nil or not properly initialized")
	}

	// Get the audio stream URL
	streamURL, err := GetAudioURL(videoURL)
	if err != nil {
		log.Printf("Error fetching audio stream URL: %v\n", err)
		return err
	}

	dgvoice.PlayAudioFile(vc, streamURL, make(chan bool))

	return nil
}
