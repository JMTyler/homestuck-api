package main

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/fcm"
	"github.com/JMTyler/homestuck-watcher/internal/scraper"
	"os"
)

// TODO: This is identical to worker.discoverNewStories() - merge them.
func populateEmptyStories() {
	for _, data := range scraper.FetchAllStories() {
		story := &db.Story{
			Domain: data["domain"],
			Endpoint: data["endpoint"],
			Collection: data["collection"],
			Title: data["title"],
			Page: 1,
		}
		story.FindOrCreate()
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No command provided")
		return
	}

	fmt.Println()
	defer fmt.Println("\n[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	switch os.Args[1] {
	case "populate":
		populateEmptyStories()
		return
	case "ping":
		endpoint := "epilogues/candy"
		if len(os.Args) >= 3 {
			endpoint = os.Args[2]
		}

		story := &db.Story{Domain: "homestuck.com", Endpoint: endpoint}
		story.Find()
		fcm.Ping(fcm.SyncEvent, story.Collection, story.Title, story.Domain, story.Endpoint, story.Page)
		return
	case "potato":
		endpoint := "epilogues/candy"
		if len(os.Args) >= 3 {
			endpoint = os.Args[2]
		}

		story := &db.Story{Domain: "homestuck.com", Endpoint: endpoint}
		story.Find()
		fcm.Ping(fcm.PotatoEvent, story.Collection, story.Title, story.Domain, story.Endpoint, story.Page)
		return
	}

	fmt.Println("Invalid command provided")
}
