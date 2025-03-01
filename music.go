package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func resolveQuery(query string) (string, error) {
	if strings.HasPrefix(query, "http://www.youtube.com") || strings.HasPrefix(query, "https://www.youtube.com") {
		youtubeRegex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/(watch\?v=|embed/|v/)?([a-zA-Z0-9_-]+)`)
		if youtubeRegex.MatchString(query) {
			return query, nil
		}
		return "", errors.New("invalid YouTube URL")
	}

	results, err := SearchYouTube(query)
	if err != nil {
		return "", err
	}

	return results[0], nil
}

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

func PlaySong(audio string, vc *discordgo.VoiceConnection, aqh *AudioQueueHandler, vh *VoiceHandler) error {
	if vc == nil || vc.ChannelID == "" {
		log.Println("Voice connection is nil or not properly initialized.")
		return fmt.Errorf("Voice connection is nil or not properly initialized")
	}

	// Remove song from queue before playing
	queue := aqh.GetQueue()
	if len(queue) > 0 {
		audio = aqh.RemoveFromQueue().URL
	}

	// Fetch stream to play
	url, err := GetAudioURL(audio)
	if err != nil {
		log.Printf("Error fetching audio stream URL: %v\n", err)
		return err
	}

	vh.connections[GUILD_ID].IsPlaying = true
	dgvoice.PlayAudioFile(vc, url, make(chan bool))
	vh.connections[GUILD_ID].IsPlaying = false

	// Need to add logic to check queue here

	return nil
}
