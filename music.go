package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jogramming/dca"
)

func PlaySong(s string, vc *discordgo.VoiceConnection) {
	err := vc.Speaking(true)
	if err != nil {
		log.Fatal("Failed setting speaking", err)
	}
	defer vc.Speaking(false)

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 120

	videoInfo, err := YTClient.GetVideo(s)
	if err != nil {
		log.Printf("Error fetching video info: %v\n", err)
	}

	// Get valid audio format
	formats := videoInfo.Formats.WithAudioChannels()
	if len(formats) == 0 {
		log.Println("No valid audio formats found.")
		return
	}

	url, err := YTClient.GetStreamURL(videoInfo, &formats[0])
	if err != nil {
		log.Printf("Error getting stream URL: %v\n", err)
		return
	}

	log.Println("Download URL:", url)

	// Encode the stream for Discord
	encodeSession, err := dca.EncodeFile(url, options)
	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}
	defer encodeSession.Cleanup()

	done := make(chan error)
	stream := dca.NewStream(encodeSession, vc, done)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occurred", err)
			}

			encodeSession.Truncate()
			return
		case <-ticker.C:
			stats := encodeSession.Stats()
			playbackPosition := stream.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r",
				playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}
