package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type TrackRequest struct {
	Url         string
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
}

type AudioQueue struct {
	Jobs chan *TrackRequest
}

func NewAudioQueue() *AudioQueue {
	return &AudioQueue{
		Jobs: make(chan *TrackRequest, 100),
	}
}

func (h *AudioQueue) StartWorkers(numWorkers int, task func(request *TrackRequest) error) {
	for i := range numWorkers {
		go func(id int) {
			for job := range h.Jobs {
				log.Printf("Worker %d is now playing: %s", id, job.Url)
				if err := task(job); err != nil {
					log.Printf("An error occurred trying to start task: %v", err)
				}
			}
		}(i)
	}
}

func (h *AudioQueue) AddToQueue(u string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.Jobs <- &TrackRequest{
		Url:         u,
		Session:     s,
		Interaction: i,
	}
}

func (h *AudioQueue) Stop() {
	close(h.Jobs)
}
