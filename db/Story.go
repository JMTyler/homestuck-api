package db

import (
	"fmt"
	"github.com/go-pg/pg/orm"
)

type Story struct {
	ID       int64
	Endpoint string `sql:", notnull, unique"`
}

func (s Story) String() string {
	return fmt.Sprintf("Story<id:%v, endpoint:'%s'>", s.ID, s.Endpoint)
}

func (s *Story) FindOrCreate() *Story {
	s.Init()

	inserted, err := DB.Model(s).Where("endpoint = ?", s.Endpoint).SelectOrInsert(s)
	// db.ModelContext(context.Context(), &models).Select()
	// db.Model(&models).SelectOrInsert()
	// res, err := db.Query(&models, "endpoint = ?", endpoint)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Query Complete. Inserted? %v  Model: %s\n", inserted, s)

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
	if DB != nil {
		return
	}

	InitDatabase()

	err := DB.CreateTable((*Story)(nil), &orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		panic(err)
	}
}
