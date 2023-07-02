package music_player

type SearchResults struct {
	results   []Track
	length    int
	maxLength int
}

func NewSearchResults() SearchResults {
	return SearchResults{
		results:   []Track{},
		length:    0,
		maxLength: 5,
	}
}

func (sr *SearchResults) AddResult(track Track) {
	if sr.length < sr.maxLength {
		sr.results = append(sr.results, track)
		sr.length++
	}
}

func (sr *SearchResults) AddResults(tracks []Track) {
	for _, track := range tracks {
		if sr.length >= sr.maxLength {
			break
		}

		sr.results = append(sr.results, track)
		sr.length++
	}
}

func (sr *SearchResults) GetResult(i int) *Track {
	if i < 0 || i >= sr.length {
		return nil
	}

	return &sr.results[i]
}

func (sr *SearchResults) GetResults() []Track {
	return sr.results
}

func (sr *SearchResults) Clear() {
	sr.results = []Track{}
	sr.length = 0
}

func (sr *SearchResults) Length() int {
	return sr.length
}

func (sr *SearchResults) MaxLength() int {
	return sr.maxLength
}
