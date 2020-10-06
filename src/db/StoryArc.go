package db

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/src/fcm"
	"github.com/go-pg/pg/orm"
	"time"
)

type StoryArc struct {
	ID        int64
	Title     string
	Endpoint  string    `pg:", notnull"`
	Page      int       `pg:", notnull"`
	StoryID   int64     `pg:", notnull, on_delete:CASCADE, on_update:CASCADE"`
	CreatedAt time.Time `pg:", notnull, default:now()"`
	UpdatedAt time.Time `pg:", notnull, default:now()"`
	Story     *Story
}

func (s StoryArc) String() string {
	return fmt.Sprintf("StoryArc<id:%v, endpoint:'%v', title:'%v', page:%v, %v>", s.ID, s.Endpoint, s.Title, s.Page, s.Story)
}

func (s *StoryArc) Scrub(version string) map[string]interface{} {
	// v1
	return map[string]interface{}{
		"endpoint": s.Endpoint,
		"title":    s.Story.Title,
		"subtitle": s.Title,
		"pages":    s.Page,
	}
}

func (a *StoryArc) FindOrCreate() *StoryArc {
	a.Init()

	_, err := DB.Model(a).Relation("Story").Where("story_arc.story_id = ? AND story_arc.endpoint = ?", a.StoryID, a.Endpoint).SelectOrInsert(a)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Query Complete. Inserted? %v  Model: %s\n", inserted, a)

	return a
}

func (a *StoryArc) Find() *StoryArc {
	a.Init()

	query := DB.Model(a).Relation("Story").Where("story_arc.endpoint = ?", a.Endpoint)
	if a.StoryID != 0 {
		query.Where("story_arc.story_id = ?", a.StoryID)
	}
	err := query.Select(a)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Query Complete. Model: %s\n", a)

	return a
}

func (a *StoryArc) Update() {
	a.Init()

	a.UpdatedAt = time.Now()

	err := DB.Update(a)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Update Complete. Model: %s\n", a)
}

func (a *StoryArc) FindAll(version string) []*StoryArc {
	a.Init()

	var arcs []*StoryArc
	query := DB.Model(&arcs).Relation("Story").Order("story_arc.created_at")
	if version == "v1" {
		query.Where("story.domain = ?", "homestuck.com")
	}
	err := query.Select()
	if err != nil {
		panic(err)
	}
	return arcs
}

func (a *StoryArc) ProcessPotato(page int) {
	// fmt.Printf("Updating story-arc #%v with Page = %v\n", a.ID, page)
	a.Page = page
	a.Update()
	fcm.Ping(fcm.PotatoEvent, a.Story.Title, a.Title, a.Story.Domain, a.Endpoint, a.Page)
}

func (a StoryArc) Init() {
	InitDatabase()

	err := DB.CreateTable((*StoryArc)(nil), &orm.CreateTableOptions{IfNotExists: true, FKConstraints: true})
	if err != nil {
		panic(err)
	}
}
