package jazzdata

import (
	"os"
	"strings"
)

type Cache struct {
	CacheFile string
	Entries   []string
	saved     bool
}

func NewCache(cacheFile string) *Cache {
	c := &Cache{
		CacheFile: cacheFile,
		saved:     true,
	}
	if cacheFile != "" {
		if _, err := os.Stat(cacheFile); err == nil {
			data, err := os.ReadFile(cacheFile)
			if err == nil {
				c.Entries = strings.Split(string(data), "\n")
			}
		}
	}
	return c
}

func (c *Cache) Add(entry string) {
	if !c.Has(entry) {
		c.saved = false
		c.Entries = append(c.Entries, entry)
	}
}

func (c *Cache) Has(entry string) bool {
	for _, v := range c.Entries {
		if v == entry {
			return true
		}
	}
	return false
}

func (c *Cache) Save() error {
	if c.saved || c.CacheFile == "" {
		return nil
	}
	data := strings.Join(c.Entries, "\n")
	err := os.WriteFile(c.CacheFile, []byte(data), 0644)
	if err == nil {
		c.saved = true
	}
	return err
}
