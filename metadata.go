// jazzdata/metadata.go
package jazzdata

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/rickcollette/jazzdata/models"
)

var SupportedSongExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".m4a":  true,
}

var (
	ErrMissingMetadata  = errors.New("missing metadata")
	ErrMetadataNotFound = errors.New("metadata not found")
)

// GetMetadataFromFile extracts detailed metadata from an audio file.
func GetMetadataFromFile(path string) (models.Metadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return models.Metadata{}, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return models.Metadata{}, err
	}

	album := m.Album()
	artist := m.Artist()
	if album == "" || artist == "" {
		return models.Metadata{}, ErrMissingMetadata
	}

	title := m.Title()
	track, totalTrack := m.Track()
	disc, totalDisc := m.Disc()

	raw := m.Raw()

	md := models.Metadata{
		Album:       album,
		Artist:      artist,
		Title:       title,
		AlbumArtist: m.AlbumArtist(),
		Composer:    m.Composer(),
		Comment:     m.Comment(),
		Genre:       m.Genre(),
		Year:        m.Year(),
		TrackNumber: track,
		TrackTotal:  totalTrack,
		DiscNumber:  disc,
		DiscTotal:   totalDisc,
		Raw:         raw,
	}

	return md, nil
}

// GetMetadataFromDirectory searches through a directory for a supported audio file
// and returns its metadata.
func GetMetadataFromDirectory(path string) (models.Metadata, error) {
	var tries int
	var result models.Metadata

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if SupportedSongExtensions[ext] {
			meta, err := GetMetadataFromFile(path)
			if err == nil {
				result = meta
				// Stop the walk once valid metadata is found.
				return filepath.SkipDir
			}
			tries++
			if tries > 3 {
				return errors.New("tries exceeded")
			}
		}
		return nil
	})
	if err != nil {
		return models.Metadata{}, err
	}
	if result.Album == "" || result.Artist == "" {
		return models.Metadata{}, ErrMetadataNotFound
	}
	return result, nil
}
