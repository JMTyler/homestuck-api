package watcher

import (
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/scraper"
)

func discoverNewStories() {
	for _, data := range scraper.FetchAllStories() {
		story := &db.Story{
			Domain:     data["domain"],
			Endpoint:   data["endpoint"],
			Collection: data["collection"],
			Title:      data["title"],
			Page:       1,
		}
		story.FindOrCreate()
	}
}

func updatePageCounts() {
	stories := new(db.Story).FindAll()
	for _, story := range stories {
		latestPage := scraper.SeekLatestPage(story.Domain, story.Endpoint, story.Page)
		if latestPage != story.Page {
			story.Potato(latestPage)
		}
	}
}
