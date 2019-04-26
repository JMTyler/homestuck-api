package db

import (
	"fmt"
	"github.com/go-pg/pg/orm"
)

type StoryArc struct {
	ID       int64
	Endpoint string `sql:", notnull, unique"`
	Page     int    `sql:", notnull"`
	StoryID  int64  `sql:", notnull, on_delete:CASCADE, on_update:CASCADE"`
	Story    *Story
}

func (s StoryArc) String() string {
	return fmt.Sprintf("StoryArc<id:%v, endpoint:'%s', page:%v, story_id:%v>", s.ID, s.Endpoint, s.Page, s.StoryID)
}

func (a *StoryArc) FindOrCreate() *StoryArc {
	a.Init()

	inserted, err := DB.Model(a).Where("endpoint = ?", a.Endpoint).SelectOrInsert(a)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Query Complete. Inserted? %v  Model: %s\n", inserted, a)

	return a
}

func (a *StoryArc) Update() {
	a.Init()

	err := DB.Update(a)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Update Complete. Model: %s\n", a)
}

func (a *StoryArc) FindAll() []StoryArc {
	a.Init()

	var arcs []StoryArc
	err := DB.Model(&arcs).Select()
	if err != nil {
		panic(err)
	}
	return arcs
}

func (a StoryArc) Init() {
	if DB != nil {
		return
	}

	InitDatabase()

	err := DB.CreateTable((*StoryArc)(nil), &orm.CreateTableOptions{IfNotExists: true, FKConstraints: true})
	if err != nil {
		panic(err)
	}
}
