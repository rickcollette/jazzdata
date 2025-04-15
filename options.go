// jazzdata/options.go
package jazzdata

import "github.com/rickcollette/jazzdata/models"

type Options struct {
	Paths              []string
	Providers          []models.Source // Now using the Source type from models.
	CoverName          string
	CachePath          string
	Tags               []string
	Recursive          bool
	Upgrade            bool
	MaxSize            float64
	MaxUpgradeSize     float64
	Strict             bool
	MaxHammingDistance int
	SilenceWarnings    bool
	DeleteOldCovers    bool
}
