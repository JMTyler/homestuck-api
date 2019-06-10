package main

import (
	"fmt"
	// que "github.com/bgentry/que-go"
	// "github.com/jackc/pgx"
	"homestuck-watcher/db"
	// "os"
	"github.com/robfig/cron"
	"homestuck-watcher/utils"
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
	c.AddFunc("0 */5 * * * *", func() {
		// TODO: once per day, do heavyweight
		// ...
		go runHeavyweightWorker()
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
