package musicplayer

import (
	"github.com/disgoorg/disgolink/lavalink"
)

type MusicQueue struct {
	tracks []lavalink.AudioTrack
	size   int
}

func NewMusicQueue() *MusicQueue {
	return &MusicQueue{
		tracks: []lavalink.AudioTrack{},
		size:   0,
	}
}

func (q *MusicQueue) Enqueue(track lavalink.AudioTrack) {
	q.tracks = append(q.tracks, track)
	q.size += 1
}

func (q *MusicQueue) EnqueueList(tracks []lavalink.AudioTrack) {
	for _, track := range tracks {
		q.tracks = append(q.tracks, track)
		q.size += 1
	}
}

func (q *MusicQueue) Dequeue() *lavalink.AudioTrack {
	return q.DequeueAtIndex(0)
}

func (q *MusicQueue) DequeueAtIndex(index int) *lavalink.AudioTrack {
	if q.Length() < index {
		return nil
	}

	track := q.tracks[index]

	if q.Length() == 1 {
		q.tracks = []lavalink.AudioTrack{}
	} else {
		q.tracks = append(q.tracks[:index], q.tracks[index+1:]...)
	}

	q.size -= 1

	return &track
}

func (q *MusicQueue) Length() int {
	return q.size
}

func (q *MusicQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *MusicQueue) Peek() *lavalink.AudioTrack {
	if q.IsEmpty() {
		return nil
	}

	return &q.tracks[0]
}

func (q *MusicQueue) PeekList() []lavalink.AudioTrack {
	return q.tracks
}
