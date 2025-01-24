package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"

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

func PlaySong(videoURL string, vc *discordgo.VoiceConnection) {
    // Validate voice connection
    if vc == nil {
        log.Println("ERROR: Voice connection is nil")
        return
    }
    if vc.ChannelID == "" {
        log.Println("ERROR: No active voice channel")
        return
    }

    // Get audio stream URL
    streamURL, err := GetAudioURL(videoURL)
    if err != nil {
        log.Printf("ERROR: Failed to get audio stream URL: %v\n", err)
        return
    }
    log.Printf("Retrieved Stream URL: %s", streamURL)

    // Prepare FFmpeg command
    cmd := exec.Command("ffmpeg",
        "-i", streamURL,
        "-f", "s16le",     // Raw audio format
        "-ar", "48000",    // Audio rate (Discord requirement)
        "-ac", "2",        // 2 audio channels
        "-acodec", "pcm_s16le",
        "pipe:1")

    // Capture command stderr for detailed error information
    var stderr bytes.Buffer
    cmd.Stderr = &stderr

    // Create stdout pipe
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Printf("ERROR: Failed to create stdout pipe: %v\n", err)
        log.Printf("STDERR: %s", stderr.String())
        return
    }

    // Start FFmpeg
    if err := cmd.Start(); err != nil {
        log.Printf("ERROR: Failed to start FFmpeg: %v\n", err)
        log.Printf("STDERR: %s", stderr.String())
        return
    }

    // Enable speaking
    if err := vc.Speaking(true); err != nil {
        log.Printf("ERROR: Failed to set speaking: %v\n", err)
        return
    }
    defer vc.Speaking(false)

    // Prepare audio buffer
    buffer := make([]byte, 960*2*2) // Adjusted buffer size

    log.Println("Starting audio streaming...")

    for {
        // Read audio data
        n, err := stdout.Read(buffer)
        if err != nil {
            if err != io.EOF {
                log.Printf("ERROR reading audio: %v\n", err)
            }
            break
        }

        if n > 0 {
            // Send audio packets
            vc.OpusSend <- buffer[:n]
        }
    }

    // Wait for FFmpeg to complete
    if err := cmd.Wait(); err != nil {
        log.Printf("ERROR: FFmpeg command failed: %v\n", err)
        log.Printf("STDERR: %s", stderr.String())
    }

    log.Println("Audio streaming completed")
}

// func PlaySonag(videoURL string, vc *discordgo.VoiceConnection) {
// 	// Get the audio stream URL
// 	streamURL, err := GetAudioURL(videoURL)
// 	if err != nil {
// 		log.Printf("Error fetching audio stream URL: %v\n", err)
// 		return
// 	}

// 	log.Println("Audio Stream URL:", streamURL)

// 	// Direct FFmpeg command to pipe audio
// 	cmd := exec.Command("ffmpeg",
// 		"-i", streamURL,
// 		"-f", "s16le",
// 		"-ar", "48000",
// 		"-ac", "2",
// 		"pipe:1")

// 	// Setup output pipe
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		log.Printf("Error creating stdout pipe: %v\n", err)
// 		return
// 	}

// 	if err := cmd.Start(); err != nil {
// 		log.Printf("Error starting FFmpeg: %v\n", err)
// 		return
// 	}

// 	// Setup Discord voice transmission
// 	if err := vc.Speaking(true); err != nil {
// 		log.Printf("Failed to set speaking: %v\n", err)
// 		return
// 	}
// 	defer vc.Speaking(false)

// 	// Send audio packets
// 	buffer := make([]byte, 960*2)
// 	for {
// 		n, err := stdout.Read(buffer)
// 		if err != nil {
// 			log.Printf("Error reading from FFmpeg: %v\n", err)
// 			break
// 		}
// 		vc.OpusSend <- buffer[:n]
// 	}

// 	cmd.Wait()
// }

// func PlaySong(videoURL string, vc *discordgo.VoiceConnection) {
// 	if vc == nil || vc.ChannelID == "" {
// 		log.Println("Voice connection is nil or not properly initialized.")
// 		return
// 	}
// 	log.Printf("Bot is connected to voice channel: %s", vc.ChannelID)

// 	err := vc.Speaking(true)
// 	if err != nil {
// 		log.Fatal("Failed setting speaking", err)
// 	}
// 	defer vc.Speaking(false)

// 	// Get the audio stream URL
// 	streamURL, err := GetAudioURL(videoURL)
// 	if err != nil {
// 		log.Printf("Error fetching audio stream URL: %v\n", err)
// 		return
// 	}

// 	log.Println("Audio Stream URL:", streamURL)

// 	// Set DCA encode options

// 	options := dca.StdEncodeOptions
// 	options.RawOutput = true
// 	options.Bitrate = 96
// 	options.Application = "lowdelay"

// 	// Encode the stream for Discord
// 	encodeSession, err := dca.EncodeFile(streamURL, options)
// 	if err != nil {
// 		log.Fatal("Failed creating an encoding session: ", err)
// 	}
// 	defer log.Println("messages", encodeSession.FFMPEGMessages())
// 	defer encodeSession.Cleanup()

// 	log.Println("Encoding session created")
// 	log.Println("encoded session", encodeSession)

// 	done := make(chan error)
// 	stream := dca.NewStream(encodeSession, vc, done)

// 	log.Println("Stream created")

// 	// ticker := time.NewTicker(time.Second)
// 	// defer ticker.Stop()

// 	ticker := time.NewTicker(500 * time.Millisecond) // Adjust as needed
// 	defer ticker.Stop()

// 	log.Println("Starting playback")
// 	log.Printf("Stats: %+v", encodeSession.Stats())
// 	log.Printf("Stream playback position: %v", stream.PlaybackPosition())

// 	for {
// 		log.Println("Streaming audio packet...")
// 		select {
// 		case err := <-done:
// 			if err != nil {
// 				log.Printf("Error received: %v", err)
// 			} else {
// 				log.Println("Playback finished successfully")
// 			}
// 			encodeSession.Cleanup()
// 			return
// 		// case err := <-done:
// 		// 	if err != nil && err != io.EOF {
// 		// 		log.Fatal("An error occurred", err)
// 		// 	}

// 		// 	log.Println("Playback finished")
// 		// 	encodeSession.Cleanup()
// 		// 	return
// 		case <-ticker.C:
// 			stats := encodeSession.Stats()
// 			playbackPosition := stream.PlaybackPosition()

// 			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r",
// 				playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
// 		}
// 	}
// }
