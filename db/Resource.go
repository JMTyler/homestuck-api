package db

import (
	"fmt"
	"github.com/go-pg/pg"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
	sql, _ := q.FormattedQuery()
	fmt.Println("[SQL]", sql)
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {}

var DB *pg.DB

func InitDatabase() *pg.DB {
	if DB != nil {
		return DB
	}

	DB = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
	})

	// DB.AddQueryHook(dbLogger{})

	return DB
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}
