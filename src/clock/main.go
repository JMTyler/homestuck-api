package main

import (
	"fmt"
	// que "github.com/bgentry/que-go"
	// "github.com/jackc/pgx"
	"homestuck-watcher/db"
	// "os"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println()
	defer fmt.Println("\n[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	c := cron.New()
	c.AddFunc("0 * * * * *", func() {
		// TODO: once every minute, do lightweight
		// ... start one-off dyno of `clock/worker lightweight`
		go runLightweightWorker()
	})
	c.AddFunc("@daily", func() {
		// TODO: once per day, do heavyweight
		// ...
	})
	c.Start()
	defer c.Stop()
	defer fmt.Println("Stopping cron ...")

	// .................

	// Catch signal so we can shutdown gracefully
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	// Wait for a signal
	<-sigCh

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
