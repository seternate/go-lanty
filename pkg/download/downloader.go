package download

import (
	"sort"

	"github.com/seternate/go-lanty-server/pkg/game"
	"golang.org/x/exp/maps"
)

type Downloader struct {
	Download map[game.Game]*Download
}

func (d *Downloader) IsDownloading(game game.Game) bool {
	download := d.Download[game]
	if download == nil {
		return false
	}
	return !download.IsComplete()
}

func (d *Downloader) Games() (keyList game.Games) {
	keys := maps.Keys(d.Download)

	sort.Slice(keys, func(i int, j int) bool {
		return keys[i].Slug < keys[j].Slug
	})
	return keys
}
