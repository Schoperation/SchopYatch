package music_player

import "github.com/disgoorg/disgolink/v2/lavalink"

type SearchResults struct {
	results   []lavalink.Track
	length    int
	maxLength int
}

func NewSearchResults() SearchResults {
	return SearchResults{
		results:   []lavalink.Track{},
		length:    0,
		maxLength: 5,
	}
}

func (sr *SearchResults) AddResult(track lavalink.Track) {
	if sr.length < sr.maxLength {
		sr.results = append(sr.results, track)
		sr.length++
	}
}

func (sr *SearchResults) AddResults(tracks []lavalink.Track) {
	for _, track := range tracks {
		if sr.length >= sr.maxLength {
			break
		}

		sr.results = append(sr.results, track)
		sr.length++
	}
}

func (sr *SearchResults) GetResult(i int) *lavalink.Track {
	if i < 0 || i >= sr.length {
		return nil
	}

	return &sr.results[i]
}

func (sr *SearchResults) GetResults() []lavalink.Track {
	return sr.results
}

func (sr *SearchResults) Clear() {
	sr.results = []lavalink.Track{}
	sr.length = 0
}

func (sr *SearchResults) Length() int {
	return sr.length
}

func (sr *SearchResults) MaxLength() int {
	return sr.maxLength
}
