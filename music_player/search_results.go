package music_player

import "github.com/disgoorg/disgolink/lavalink"

type SearchResults struct {
	results   []lavalink.AudioTrack
	maxLength int
}

func NewSearchResults() SearchResults {
	return SearchResults{
		results:   []lavalink.AudioTrack{},
		maxLength: 5,
	}
}

func (sr *SearchResults) AddTrack(track lavalink.AudioTrack) {
	if sr.Length() < sr.MaxLength() {
		sr.results = append(sr.results, track)
	}
}

func (sr *SearchResults) GetTrack(i int) *lavalink.AudioTrack {
	if i < 0 || i >= sr.Length() {
		return nil
	}

	return &sr.results[i]
}

func (sr *SearchResults) Clear() {
	sr.results = []lavalink.AudioTrack{}
}

func (sr *SearchResults) Length() int {
	return len(sr.results)
}

func (sr *SearchResults) MaxLength() int {
	return sr.maxLength
}
