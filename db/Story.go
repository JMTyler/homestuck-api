package db

import (
	"fmt"
	"github.com/go-pg/pg/orm"
	"time"
)

type Story struct {
	ID        int64
	Title     string
	Endpoint  string    `sql:", notnull, unique"`
	CreatedAt time.Time `sql:", notnull, default:now()"`
	UpdatedAt time.Time `sql:", notnull, default:now()"`
}

func (s Story) String() string {
	return fmt.Sprintf("Story<id:%v, endpoint:'%s', title:'%s'>", s.ID, s.Endpoint, s.Title)
}

func (s *Story) FindOrCreate() *Story {
	s.Init()

	_, err := DB.Model(s).Where("endpoint = ?", s.Endpoint).SelectOrInsert(s)
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

func (s Story) Init() {
	InitDatabase()

	err := DB.CreateTable((*Story)(nil), &orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		panic(err)
	}
}
