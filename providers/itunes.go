// jazzdata/providers/itunes.go
package providers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rickcollette/jazzdata/models"
)

type ITunesProvider struct{}

func (ip ITunesProvider) GetCovers(artist, album string) ([]models.Cover, error) {
	baseURL := "https://itunes.apple.com/search"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("term", artist)
	q.Add("media", "music")
	q.Add("entity", "album")
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data struct {
		Results []struct {
			CollectionName string `json:"collectionName"`
			ArtworkUrl100  string `json:"artworkUrl100"`
			ArtworkUrl60   string `json:"artworkUrl60"`
			ArtistName     string `json:"artistName"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	var covers []models.Cover
	for _, item := range data.Results {
		if strings.Contains(strings.ToLower(item.CollectionName), strings.ToLower(album)) {
			coverURL := item.ArtworkUrl100
			if coverURL == "" {
				coverURL = item.ArtworkUrl60
			}
			cover := models.Cover{
				Artist:     item.ArtistName,
				Title:      item.CollectionName,
				Source:     models.ITunes,
				CoverURL:   coverURL,
				Ext:        "", // Optionally, compute extension.
				Confidence: 100,
			}
			covers = append(covers, cover)
		}
	}
	return covers, nil
}
