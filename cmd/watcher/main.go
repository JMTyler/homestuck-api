package main

import (
	// que "github.com/bgentry/que-go"
	// "github.com/jackc/pgx"
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/utils"
	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println()
	defer fmt.Println("\n[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	// TODO: Consider configuring these from a database table.
	c := cron.New(cron.WithSeconds())
	c.AddFunc("5 * * * * *", func() {
		// TODO: convert to one-off dyno of `clock/worker lightweight`
		go updatePageCounts()
	})
	c.AddFunc("0 0 * * * *", func() {
		go discoverNewStories()
	})
	c.Start()
	defer c.Stop()
	defer fmt.Println("Stopping cron ...")

	// .................

	utils.GracefulShutdown()

	// pgxcfg, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	panic(err)
	// }

	// pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
	// 	ConnConfig:   pgxcfg,
	// 	AfterConnect: que.PrepareStatements,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer pgxpool.Close()

	// qc := que.NewClient(pgxpool)
	// qc.Enqueue(&que.Job{
	// 	Type: "Lightweight",
	// })
}
