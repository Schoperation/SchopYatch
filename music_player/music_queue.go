package music_player

import (
	"math/rand"
	"time"

	"github.com/disgoorg/disgolink/v2/lavalink"
)

type MusicQueue struct {
	tracks   []lavalink.Track
	size     int
	duration lavalink.Duration
}

func NewMusicQueue() MusicQueue {
	return MusicQueue{
		tracks:   []lavalink.Track{},
		size:     0,
		duration: lavalink.Duration(0),
	}
}

func (q *MusicQueue) Enqueue(track lavalink.Track) {
	q.tracks = append(q.tracks, track)
	q.size++
	q.duration += track.Info.Length
}

func (q *MusicQueue) EnqueueList(tracks []lavalink.Track) {
	for _, track := range tracks {
		q.tracks = append(q.tracks, track)
		q.size++
		q.duration += track.Info.Length
	}
}

func (q *MusicQueue) Dequeue() *lavalink.Track {
	return q.DequeueAt(0)
}

func (q *MusicQueue) DequeueAt(index int) *lavalink.Track {
	if q.Length() <= index {
		return nil
	}

	track := q.tracks[index]

	if q.Length() == 1 {
		q.tracks = []lavalink.Track{}
	} else {
		q.tracks = append(q.tracks[:index], q.tracks[index+1:]...)
	}

	q.size--
	q.duration -= track.Info.Length

	return &track
}

func (q *MusicQueue) Length() int {
	return q.size
}

func (q *MusicQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *MusicQueue) Duration() lavalink.Duration {
	return q.duration
}

func (q *MusicQueue) Peek() *lavalink.Track {
	return q.PeekAt(0)
}

func (q *MusicQueue) PeekAt(index int) *lavalink.Track {
	if q.Length() <= index {
		return nil
	}

	return &q.tracks[index]
}

func (q *MusicQueue) PeekList() []lavalink.Track {
	return q.tracks
}

func (q *MusicQueue) Clear() {
	q.tracks = []lavalink.Track{}
	q.size = 0
	q.duration = lavalink.Duration(0)
}

func (q *MusicQueue) Shuffle() {
	if q.IsEmpty() {
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixMicro()))
	rng.Shuffle(q.Length(), func(i, j int) {
		q.tracks[i], q.tracks[j] = q.tracks[j], q.tracks[i]
	})
}
