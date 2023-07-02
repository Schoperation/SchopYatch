package music_player

import (
	"math/rand"
	"time"
)

type MusicQueue struct {
	tracks   []Track
	size     int
	duration time.Duration
}

func NewMusicQueue() MusicQueue {
	return MusicQueue{
		tracks:   []Track{},
		size:     0,
		duration: time.Duration(0),
	}
}

func (q *MusicQueue) Enqueue(track Track) {
	q.tracks = append(q.tracks, track)
	q.size++
	q.duration += track.Length()
}

func (q *MusicQueue) EnqueueList(tracks []Track) {
	for _, track := range tracks {
		q.tracks = append(q.tracks, track)
		q.size++
		q.duration += track.Length()
	}
}

func (q *MusicQueue) Dequeue() *Track {
	return q.DequeueAt(0)
}

func (q *MusicQueue) DequeueAt(index int) *Track {
	if q.Length() <= index {
		return nil
	}

	track := q.tracks[index]

	if q.Length() == 1 {
		q.tracks = []Track{}
	} else {
		q.tracks = append(q.tracks[:index], q.tracks[index+1:]...)
	}

	q.size--
	q.duration -= track.Length()

	return &track
}

func (q *MusicQueue) Length() int {
	return q.size
}

func (q *MusicQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *MusicQueue) Duration() time.Duration {
	return q.duration
}

func (q *MusicQueue) Peek() *Track {
	return q.PeekAt(0)
}

func (q *MusicQueue) PeekAt(index int) *Track {
	if q.Length() <= index {
		return nil
	}

	return &q.tracks[index]
}

func (q *MusicQueue) PeekList() []Track {
	return q.tracks
}

func (q *MusicQueue) Clear() {
	q.tracks = []Track{}
	q.size = 0
	q.duration = time.Duration(0)
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
