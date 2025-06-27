package services

import (
	"bytes"
	"os/exec"
	"strings"
)

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
