package music_player

import (
	"strings"
	"time"

	"github.com/disgoorg/disgolink/v3/lavalink"
)

type Track struct {
	title  string
	author string
	length time.Duration

	encoded    string
	identifier string
	isStream   bool
	uri        *string
	sourceName string
	position   time.Duration
}

// NewTrack mostly for testing purposes
func NewTrack(encoded, title, author string, length time.Duration) Track {
	return Track{
		encoded: encoded,
		title:   title,
		author:  author,
		length:  length,
	}
}

func (track Track) Title() string {
	title := strings.ReplaceAll(track.title, "\\", "\\\\")
	title = strings.ReplaceAll(title, "*", "\\*")
	title = strings.ReplaceAll(title, "_", "\\_")
	return title
}

func (track Track) Author() string {
	author := strings.ReplaceAll(track.author, "\\", "\\\\")
	author = strings.ReplaceAll(author, "*", "\\*")
	author = strings.ReplaceAll(author, "_", "\\_")
	return author
}

func (track Track) Length() time.Duration {
	return track.length
}

func (track Track) ToLavalinkTrack() lavalink.Track {
	return lavalink.Track{
		Encoded: track.encoded,
		Info: lavalink.TrackInfo{
			Identifier: track.identifier,
			Author:     track.author,
			Length:     toLavalinkDuration(track.length),
			IsStream:   track.isStream,
			Title:      track.title,
			URI:        track.uri,
			SourceName: track.sourceName,
			Position:   toLavalinkDuration(track.position),
		},
	}
}

func toTrack(track lavalink.Track) Track {
	return Track{
		encoded:    track.Encoded,
		identifier: track.Info.Identifier,
		author:     track.Info.Author,
		length:     toTimeDuration(track.Info.Length),
		isStream:   track.Info.IsStream,
		title:      track.Info.Title,
		uri:        track.Info.URI,
		sourceName: track.Info.SourceName,
		position:   toTimeDuration(track.Info.Position),
	}
}

func toTracks(lavalinkTracks []lavalink.Track) []Track {
	tracks := make([]Track, len(lavalinkTracks))
	for i := range lavalinkTracks {
		tracks[i] = toTrack(lavalinkTracks[i])
	}

	return tracks
}

func toLavalinkDuration(duration time.Duration) lavalink.Duration {
	return lavalink.Duration(duration.Milliseconds())
}

func toTimeDuration(duration lavalink.Duration) time.Duration {
	return time.Duration(duration.Milliseconds() * 1000000)
}
