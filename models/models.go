// jazzdata/models/models.go
package models

// Cover represents cover art data.
type Cover struct {
	Artist     string
	Title      string
	Source     Source
	CoverURL   string
	Ext        string
	Confidence int
}

// Source represents the music provider (as a string).
type Source string

const (
	AppleMusic Source = "applemusic"
	Deezer     Source = "deezer"
	Discogs    Source = "discogs"
	ITunes     Source = "itunes"
)

// Metadata holds album metadata.
type Metadata struct {
	Album  string
	Artist string
}
