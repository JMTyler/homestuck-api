package db

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/fcm"
	"github.com/go-pg/pg/orm"
	"time"
)

type Story struct {
	ID         int64     `json:"-"`
	Collection string    `json:"title"`
	Title      string    `json:"subtitle"`
	Domain     string    `json:"-"        pg:", notnull, unique:ix_domain_endpoint"`
	Endpoint   string    `json:"endpoint" pg:", notnull, unique:ix_domain_endpoint"`
	Page       int       `json:"pages"    pg:", notnull"`
	CreatedAt  time.Time `json:"-"        pg:", notnull, default:now()"`
	UpdatedAt  time.Time `json:"-"        pg:", notnull, default:now()"`
}

func (s *Story) String() string {
	title := s.Title
	if s.Collection != "" {
		title = s.Collection + ": " + title
	}
	return fmt.Sprintf("Story<url:'%s', title:'%s'>", s.Domain+"/"+s.Endpoint, title)
}

func (s *Story) FindOrCreate() *Story {
	s.Init()

	_, err := DB.Model(s).Where("domain = ? AND endpoint = ?", s.Domain, s.Endpoint).SelectOrInsert(s)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *Story) Find() *Story {
	s.Init()

	err := DB.Model(s).Where("domain = ? AND endpoint = ?", s.Domain, s.Endpoint).Select(s)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *Story) Update() {
	s.Init()

	s.UpdatedAt = time.Now()

	err := DB.Update(s)
	if err != nil {
		panic(err)
	}
}

func (s *Story) FindAll() []*Story {
	s.Init()

	var stories []*Story
	err := DB.Model(&stories).Order("created_at").Where("domain = 'homestuck.com'").Select()
	if err != nil {
		panic(err)
	}
	return stories
}

func (s *Story) Potato(page int) {
	s.Page = page
	s.Update()
	fcm.Ping(fcm.PotatoEvent, s.Collection, s.Title, s.Domain, s.Endpoint, s.Page)
}

func (s *Story) Init() {
	InitDatabase()

	err := DB.CreateTable((*Story)(nil), &orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		panic(err)
	}
}
