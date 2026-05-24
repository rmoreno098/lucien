package bot

type AudioQueue struct {
	CommandsPool *WorkerPool
	MusicPool    *WorkerPool
}

func NewAudioQueue() *AudioQueue {
	cmd := NewCommandsPool()
	music := NewMusicPool()

	cmd.Start()
	music.Start()

	return &AudioQueue{
		CommandsPool: cmd,
		MusicPool:    music,
	}
	log.Printf("Could not find voice connection for guild: %v", guildID)
	return nil
}
