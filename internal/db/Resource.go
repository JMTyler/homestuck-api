package db

import (
	"fmt"
	"github.com/go-pg/pg"
	"net/url"
	"os"
	"strings"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
	sql, _ := q.FormattedQuery()
	fmt.Println("[SQL]", sql)
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {}

var DB *pg.DB

type Model interface {
	Scrub() map[string]interface{}
}

// TODO: Cannot pass in slice of typed models, as slices cannot be converted. More reason to use collections and/or switch to json tags.
func ScrubModels(models []Model) []map[string]interface{} {
	scrubbed := make([]map[string]interface{}, len(models))
	for i, model := range models {
		scrubbed[i] = model.Scrub()
	}
	return scrubbed
}

func InitDatabase() *pg.DB {
	if DB != nil {
		return DB
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		// TODO: throw error, need database
	}

	auth, err := url.Parse(databaseUrl)
	if err != nil {
		panic(err)
	}

	if auth.Scheme != "postgres" {
		// TODO: throw error, DB must be postgres
	}

	password, _ := auth.User.Password()
	database := strings.TrimPrefix(auth.Path, "/")

	connOptions := &pg.Options{
		User:     auth.User.Username(),
		Password: password,
		Addr:     auth.Host,
		Database: database,
		// ApplicationName: "HomestuckAPI",
	}

	DB = pg.Connect(connOptions)

	// DB.AddQueryHook(dbLogger{})

	return DB
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}
