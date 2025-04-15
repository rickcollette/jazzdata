// jazzdata/providers/discogs.go
package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rickcollette/jazzdata/models"
)

type DiscogsProvider struct{}

const (
	discogsBaseURL = "https://api.discogs.com"
	discogsAPIKey  = "yoVukqDuMTyrckrqJdfc"
	discogsSecret  = "PRCvdLuRMVghFrNtRRvylkDZEKCiLUbI"
)

func (dp DiscogsProvider) GetCovers(artist, album string) ([]models.Cover, error) {
	url := discogsBaseURL + "/database/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "jazzdata/1.0")
	auth := fmt.Sprintf("Discogs key=%s, secret=%s", discogsAPIKey, discogsSecret)
	req.Header.Set("Authorization", auth)
	q := req.URL.Query()
	q.Add("artist", artist)
	q.Add("release_title", album)
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
			Title      string `json:"title"`
			CoverImage string `json:"cover_image"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	var covers []models.Cover
	target := strings.ToLower(fmt.Sprintf("%s - %s", artist, album))
	for _, item := range data.Results {
		if strings.Contains(strings.ToLower(item.Title), target) {
			cover := models.Cover{
				Artist:     artist,
				Title:      album,
				Source:     models.Discogs,
				CoverURL:   item.CoverImage,
				Ext:        "",  // Compute extension if desired.
				Confidence: 100, // Dummy value.
			}
			covers = append(covers, cover)
		}
	}
	return covers, nil
}
