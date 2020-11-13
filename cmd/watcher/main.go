package main

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/utils"
	"github.com/JMTyler/homestuck-watcher/internal/watcher"
)

func main() {
	defer db.CloseDatabase()

	watcher.Start()
	defer watcher.Stop()

	fmt.Println("Watcher has started")

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
