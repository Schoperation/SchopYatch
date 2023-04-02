package music_player

import (
	"log"

	"github.com/disgoorg/disgolink/lavalink"
)

type EventListener struct {
	Queue           *MusicQueue
	LoopMode        *LoopMode
	GotDisconnected *bool
}

func NewEventListener(queue *MusicQueue, loopMode *LoopMode, gotDisconnected *bool) lavalink.PlayerEventListener {
	return &EventListener{
		Queue:           queue,
		LoopMode:        loopMode,
		GotDisconnected: gotDisconnected,
	}
}

func (l *EventListener) OnPlayerPause(player lavalink.Player) {
	//log.Printf("OnPlayerPause")
}
func (l *EventListener) OnPlayerResume(player lavalink.Player) {
	//log.Printf("OnPlayerResume")
}
func (l *EventListener) OnPlayerUpdate(player lavalink.Player, state lavalink.PlayerState) {
	//log.Printf("OnPlayerUpdate")
}
func (l *EventListener) OnTrackStart(player lavalink.Player, track lavalink.AudioTrack) {
	//log.Printf("%d", *l.LoopMode)
}
func (l *EventListener) OnTrackEnd(player lavalink.Player, track lavalink.AudioTrack, endReason lavalink.AudioTrackEndReason) {
	if !endReason.MayStartNext() || (l.Queue.IsEmpty() && *l.LoopMode != LoopTrack) {
		return
	}

	var nextTrack lavalink.AudioTrack
	switch *l.LoopMode {
	case LoopTrack:
		nextTrack = track
	case LoopQueue:
		l.Queue.Enqueue(track)
		nextTrack = *l.Queue.Dequeue()
	default:
		nextTrack = *l.Queue.Dequeue()
	}

	err := player.Play(nextTrack)
	if err != nil {
		log.Printf("Could not play the next track: %v", err)
		return
	}
}
func (l *EventListener) OnTrackException(player lavalink.Player, track lavalink.AudioTrack, exception lavalink.FriendlyException) {
	//log.Printf("OnTrackException")
}
func (l *EventListener) OnTrackStuck(player lavalink.Player, track lavalink.AudioTrack, thresholdMs lavalink.Duration) {
	//log.Printf("OnTrackStuck")
}
func (l *EventListener) OnWebSocketClosed(player lavalink.Player, code int, reason string, byRemote bool) {
	//log.Printf("OnWebSocketClosed code: %v\n", code)
	*l.GotDisconnected = true
}
