package main

import (
	"log"
	"sync"
)

type AudioTrack struct {
	URL string
}

type AudioQueueHandler struct {
	mu    sync.Mutex
	queue []*AudioTrack
}

func NewAudioQueueHandler() *AudioQueueHandler {
	log.Println("Queue initalized.")
	return &AudioQueueHandler{
		queue: []*AudioTrack{},
	}
}

func (aqh *AudioQueueHandler) AddToQueue(url string) string {
	aqh.mu.Lock()
	defer aqh.mu.Unlock()

	track := &AudioTrack{URL: url}
	aqh.queue = append(aqh.queue, track)
	log.Printf("Added %s to queue.\n", url)
	return "Added " + url + " to the queue."
}

func (aqh *AudioQueueHandler) ClearQueue() {
	aqh.mu.Lock()
	defer aqh.mu.Unlock()

	clear(aqh.queue)
	log.Println("Queue cleared.")
}

func (aqh *AudioQueueHandler) RemoveFromQueue() *AudioTrack {
	aqh.mu.Lock()
	defer aqh.mu.Unlock()

	if len(aqh.queue) == 0 {
		log.Println("Queue is empty.")
		return nil
	}

	popped := aqh.queue[0]
	aqh.queue = aqh.queue[1:]
	log.Printf("Popped value from queue: %s.", popped.URL)
	return popped
}

func (aqh *AudioQueueHandler) GetQueue() []*AudioTrack {
	aqh.mu.Lock()
	defer aqh.mu.Unlock()

	log.Println("Successfully returned queue.")
	return aqh.queue
}
