// jazzdata/providers/provider.go
package providers

import "github.com/rickcollette/jazzdata/models"

type Provider interface {
	GetCovers(artist, album string) ([]models.Cover, error)
}
