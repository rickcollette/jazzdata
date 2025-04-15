// jazzdata/providers/apple_music.go
package providers

import (
	"net/http"
	"strings"

	"github.com/rickcollette/jazzdata/models"
)

// Export the field Itunes so it can be set from outside the package.
type AppleMusicProvider struct {
	Itunes ITunesProvider
}

func (amp AppleMusicProvider) transformURL(itunesURL string) string {
	coverURL := strings.Replace(itunesURL, "is1-ssl", "a1", -1)
	coverURL = strings.Replace(coverURL, "/image/", "/us/", -1)
	coverURL = strings.Replace(coverURL, "/thumb/", "/r1000/063/", -1)
	parts := strings.Split(coverURL, "/")
	if len(parts) > 1 {
		coverURL = strings.Join(parts[:len(parts)-1], "/")
	}
	return coverURL
}

func testURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func (amp AppleMusicProvider) GetCovers(artist, album string) ([]models.Cover, error) {
	covers, err := amp.Itunes.GetCovers(artist, album)
	if err != nil {
		return nil, err
	}
	var results []models.Cover
	for _, cover := range covers {
		tURL := amp.transformURL(cover.CoverURL)
		if testURL(tURL) {
			cover.CoverURL = tURL
			cover.Ext = "" // Optionally, determine the extension.
			cover.Source = models.AppleMusic
			results = append(results, cover)
		}
	}
	return results, nil
}
