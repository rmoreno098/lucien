package bot

import (
	"github.com/bwmarrin/discordgo"
)

type AudioQueue struct {
	CommandsPool *Commands
	MusicPool    *Music
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
}

func (q *Music) AddTask(u string, s *discordgo.Session, i *discordgo.InteractionCreate, v *VoiceHandler) {
	q.Tasks <- &PlaySong{
		Url:         u,
		Session:     s,
		Interaction: i,
		Voice:       v,
	}

	if len(q.MusicMap) > 0 {
		q.MusicMap = q.MusicMap[1:]
	}

	q.MusicMap = append(q.MusicMap, u)
}

func (q *Commands) EndQueue(s *discordgo.Session, i *discordgo.InteractionCreate, v *VoiceHandler) {
	q.Tasks <- &StopQueue{
		Session:     s,
		Interaction: i,
		Voice:       v,
	}
}
