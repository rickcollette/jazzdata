// jazzdata/providers/deezer.go
package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rickcollette/jazzdata/models"
)

type DeezerProvider struct{}

func (dp DeezerProvider) GetCovers(artist, album string) ([]models.Cover, error) {
	baseURL := "https://api.deezer.com"
	query := fmt.Sprintf(`artist:"%s" album:"%s"`, artist, album)
	url := fmt.Sprintf("%s/search/album?q=%s", baseURL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
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
		Data []struct {
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
			Title   string `json:"title"`
			CoverXL string `json:"cover_xl"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	var covers []models.Cover
	for _, item := range data.Data {
		cover := models.Cover{
			Artist:     item.Artist.Name,
			Title:      item.Title,
			Source:     models.Deezer,
			CoverURL:   item.CoverXL,
			Ext:        "", // (Optional) determine extension from URL.
			Confidence: 1,
		}
		covers = append(covers, cover)
	}
	return covers, nil
}
