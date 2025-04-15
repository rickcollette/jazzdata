// cli.go (in the module root, package main)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rickcollette/jazzdata"
	"github.com/rickcollette/jazzdata/models"
	"github.com/rickcollette/jazzdata/providers"
)

func main() {
	var providerStr string
	var coverName string
	var cachePath string
	var tagStr string
	var recursive bool
	var upgrade bool
	var maxSize float64
	var maxUpgradeSize float64
	var strict bool
	var maxHammingDistance int
	var silenceWarnings bool
	var deleteOldCovers bool

	flag.StringVar(&providerStr, "provider", "applemusic,deezer", "Comma-separated list of providers.")
	flag.StringVar(&coverName, "cover-name", "cover", "Filename for the cover art image.")
	flag.StringVar(&cachePath, "cache", "", "Path to cache file.")
	flag.StringVar(&tagStr, "tag", "", "Comma-separated list of tags.")
	flag.BoolVar(&recursive, "recursive", false, "Enable recursive album search.")
	flag.BoolVar(&upgrade, "upgrade", false, "Upgrade existing cover art.")
	flag.Float64Var(&maxSize, "max-size", 10, "Maximum size (in MB) for existing cover art.")
	flag.Float64Var(&maxUpgradeSize, "max-upgrade-size", 15, "Maximum candidate size (in MB) for upgrade.")
	flag.BoolVar(&strict, "strict", false, "Enable strict mode for upgrades.")
	flag.IntVar(&maxHammingDistance, "max-hamming-distance", 4, "Maximum allowed hamming distance.")
	flag.BoolVar(&silenceWarnings, "silence-warnings", false, "Silence warnings.")
	flag.BoolVar(&deleteOldCovers, "delete-old-covers", false, "Delete old covers (instead of renaming).")
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		fmt.Println("Please specify at least one path.")
		os.Exit(1)
	}

	// Parse providers.
	var selectedProviders []providers.Provider
	for _, p := range strings.Split(providerStr, ",") {
		p = strings.TrimSpace(p)
		switch p {
		case "applemusic":
			// Use the exported field Itunes.
			selectedProviders = append(selectedProviders, providers.AppleMusicProvider{Itunes: providers.ITunesProvider{}})
		case "deezer":
			selectedProviders = append(selectedProviders, providers.DeezerProvider{})
		case "discogs":
			selectedProviders = append(selectedProviders, providers.DiscogsProvider{})
		case "itunes":
			selectedProviders = append(selectedProviders, providers.ITunesProvider{})
		}
	}

	opts := jazzdata.Options{
		Paths:              paths,
		Providers:          []models.Source{}, // This field isn’t used directly in this sample.
		CoverName:          coverName,
		CachePath:          cachePath,
		Tags:               strings.Split(tagStr, ","),
		Recursive:          recursive,
		Upgrade:            upgrade,
		MaxSize:            maxSize,
		MaxUpgradeSize:     maxUpgradeSize,
		Strict:             strict,
		MaxHammingDistance: maxHammingDistance,
		SilenceWarnings:    silenceWarnings,
		DeleteOldCovers:    deleteOldCovers,
	}

	var pathLocations []string
	if recursive {
		albumPaths, err := jazzdata.GetAlbumPaths(paths[0], upgrade)
		if err != nil {
			fmt.Println("Error scanning directory:", err)
			os.Exit(1)
		}
		pathLocations = albumPaths
	} else {
		pathLocations = paths
	}

	if upgrade {
		// For brevity, the upgrade functionality is not fully implemented here.
		fmt.Println("Upgrade functionality not fully implemented in this sample.")
	} else {
		handleDownload(opts, pathLocations, selectedProviders)
	}
}

func handleDownload(opts jazzdata.Options, pathLocations []string, selectedProviders []providers.Provider) {
	completed := 0
	failed := 0

	for _, path := range pathLocations {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("Error reading path:", err)
			continue
		}

		var meta models.Metadata
		if info.IsDir() {
			fmt.Printf("Fetching metadata from album directory %s\n", path)
			meta, err = jazzdata.GetMetadataFromDirectory(path)
		} else {
			fmt.Printf("Fetching metadata from track %s\n", path)
			meta, err = jazzdata.GetMetadataFromFile(path)
		}
		if err != nil {
			fmt.Printf("Error: %v for %s\n", err, path)
			failed++
			continue
		}

		if jazzdata.HasCover(path) {
			fmt.Printf("Warning: %s already has a cover. Skipping.\n", path)
			continue
		}

		var results []models.Cover
		for _, prov := range selectedProviders {
			covers, err := prov.GetCovers(meta.Artist, meta.Album)
			if err == nil && len(covers) > 0 {
				results = covers
				break
			}
		}
		if len(results) == 0 {
			fmt.Printf("Error: No suitable cover found for %s\n", path)
			failed++
			continue
		}

		// Select the first valid cover.
		var cover models.Cover
		for _, c := range results {
			if _, ok := jazzdata.ImageExtensions[strings.ToLower(c.Ext)]; ok {
				cover = c
				break
			}
		}
		if cover.CoverURL == "" {
			fmt.Printf("Error: No valid cover art for %s\n", path)
			failed++
			continue
		}

		var targetDir string
		if info.IsDir() {
			targetDir = path
		} else {
			targetDir = filepath.Dir(path)
		}

		err = jazzdata.DownloadCover(cover.CoverURL, targetDir, opts.CoverName+cover.Ext)
		if err != nil {
			fmt.Printf("Error downloading cover for %s: %v\n", path, err)
			failed++
			continue
		}
		fmt.Printf("Successfully downloaded cover art for %s\n", path)
		completed++
	}
	fmt.Printf("Completed: %d, Failed: %d\n", completed, failed)
}
