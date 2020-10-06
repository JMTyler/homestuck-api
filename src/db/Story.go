package db

import (
	"fmt"
	"github.com/go-pg/pg/orm"
	"time"
)

type Story struct {
	ID         int64
	Title      string
	Domain     string    `pg:", notnull"`
	Endpoint   string    `pg:", notnull"`
	CreatedAt  time.Time `pg:", notnull, default:now()"`
	UpdatedAt  time.Time `pg:", notnull, default:now()"`
}

func (s *Story) String() string {
	return fmt.Sprintf("Story<id:%v, url:'%s', title:'%s'>", s.ID, s.Domain+"/"+s.Endpoint, s.Title)
}

func (s *Story) FindOrCreate() *Story {
	s.Init()

	_, err := DB.Model(s).Where("domain = ? AND endpoint = ?", s.Domain, s.Endpoint).SelectOrInsert(s)
	// db.ModelContext(context.Context(), &models).Select()
	// db.Model(&models).SelectOrInsert()
	// res, err := db.Query(&models, "endpoint = ?", endpoint)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Query Complete. Inserted? %v  Model: %s\n", inserted, s)

	// if res.RowsReturned() == 0 {
	// 	model = &Story{Endpoint: endpoint}
	// 	err := db.Insert(model)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// fmt.Printf("Finished Model: %s\n", model)

	return s
}

func (s *Story) Init() {
	InitDatabase()

	err := DB.CreateTable((*Story)(nil), &orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		panic(err)
	}
}
