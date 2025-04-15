// jazzdata/models/models.go
package models

// Cover represents cover art data.
type Cover struct {
	Artist     string  `json:"artist"`
	Title      string  `json:"title"`
	Source     Source  `json:"source"`
	CoverURL   string  `json:"coverURL"`
	Ext        string  `json:"ext"`
	Confidence int     `json:"confidence"`
}

// Source represents the music provider as a string.
type Source string

const (
	AppleMusic Source = "applemusic"
	Deezer     Source = "deezer"
	Discogs    Source = "discogs"
	ITunes     Source = "itunes"
)

// Metadata holds detailed album metadata.
type Metadata struct {
	Album       string                 `json:"album"`
	Artist      string                 `json:"artist"`
	Title       string                 `json:"title"`
	AlbumArtist string                 `json:"albumArtist,omitempty"`
	Composer    string                 `json:"composer,omitempty"`
	Comment     string                 `json:"comment,omitempty"`
	Genre       string                 `json:"genre,omitempty"`
	Year        int                    `json:"year,omitempty"`
	TrackNumber int                    `json:"trackNumber,omitempty"`
	TrackTotal  int                    `json:"trackTotal,omitempty"`
	DiscNumber  int                    `json:"discNumber,omitempty"`
	DiscTotal   int                    `json:"discTotal,omitempty"`
	Raw         map[string]interface{} `json:"raw,omitempty"`
}
