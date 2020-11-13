package watcher

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

var c = cron.New(cron.WithSeconds())

func Start() {
	// TODO: Consider configuring these from a database table.
	// Five seconds after every minute.
	c.AddFunc("5 * * * * *", func() {
		// TODO: convert to one-off dyno of `clock/worker lightweight`
		fmt.Println("* Updating page counts *")
		go updatePageCounts()
	})

	// Every hour.
	c.AddFunc("0 0 * * * *", func() {
		fmt.Println("* Discovering new stories *")
		go discoverNewStories()
	})

	c.Start()
}

func Stop() {
	fmt.Println("Stopping cron ...")
	c.Stop()
}
