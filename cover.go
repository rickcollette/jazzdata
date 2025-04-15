package jazzdata

type Cover struct {
	Artist     string
	Title      string
	Source     string // e.g. "applemusic", "deezer", etc.
	CoverURL   string
	Ext        string
	Confidence int
}
