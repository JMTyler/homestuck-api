package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	// "io/ioutil"
	"net/http"
	// "time"
	"regexp"
)

const BaseURL = "https://www.homestuck.com"

type Story struct {
	ID       int64
	Endpoint string `sql:", notnull, unique"`
}

func (s Story) String() string {
	return fmt.Sprintf("Story<id:%v, endpoint:'%s'>", s.ID, s.Endpoint)
}

type StoryArc struct {
	ID       int64
	Endpoint string `sql:", notnull, unique"`
	Page     int    `sql:", notnull"`
	StoryID  int64  `sql:", notnull, on_delete:CASCADE, on_update:CASCADE"`
	Story    *Story
}

func (s StoryArc) String() string {
	return fmt.Sprintf("StoryArc<id:%v, endpoint:'%s', page:%v, story_id:%v", s.ID, s.Endpoint, s.Page, s.StoryID)
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
	sql, _ := q.FormattedQuery()
	fmt.Println("[SQL]", sql)
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {}

func prepareDatabase() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
	})

	err := db.CreateTable((*Story)(nil), &orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		panic(err)
	}

	err = db.CreateTable((*StoryArc)(nil), &orm.CreateTableOptions{IfNotExists: true, FKConstraints: true})
	if err != nil {
		panic(err)
	}

	db.AddQueryHook(dbLogger{})

	return db
}

func fetch(endpoint string) *goquery.Document {
	response, err := http.Get(BaseURL + endpoint)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Status Code:", response.StatusCode)
	// fmt.Println("Body:", response.Body)

	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	panic(err)
	// }

	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		panic(err)
	}

	return doc
}

func lookupStories() []string {
	doc := fetch("/stories")

	links := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if !exists {
			return false
		}
		matched, _ := regexp.MatchString("^/log/", href)
		if !matched {
			return false
		}
		return true
	}).Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Attr("href")
		// html, _ := s.Html()
		// fmt.Printf("LINK:  %s  --  %s\n", html, href)
		return regexp.MustCompile("^/log").ReplaceAllString(href, "")
	})

	return links
}

func lookupStoryArcs(endpoint string) []string {
	doc := fetch("/log" + endpoint)

	links := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if !exists {
			return false
		}
		matched, _ := regexp.MatchString("/\\d+$", href)
		if !matched {
			return false
		}
		return true
	}).Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Attr("href")
		// html, _ := s.Html()
		// fmt.Printf("LINK:  %s  --  %s\n", html, href)
		return regexp.MustCompile("/\\d+$").ReplaceAllString(href, "")
	})

	return uniq(links)
}

func lookupLatestPage(endpoint string, page int) int {
	response, err := http.Head(fmt.Sprintf("%s%s/%v", BaseURL, endpoint, page))
	fmt.Printf("Seeking on %s -- Checking page %v -- Status Code %v\n", endpoint, page, response.StatusCode)
	if err != nil {
		panic(err)
	}

	if page == 150 {
		panic("Blew past the actual latest page!")
	}

	if response.StatusCode == 404 {
		return page - 1
	}

	if response.StatusCode == 200 {
		return lookupLatestPage(endpoint, page+1)
	}

	panic(fmt.Sprintf("Request to %s/%v returned unexpected Status Code %v\n", endpoint, page, response.StatusCode))
}

func ensureStory(db *pg.DB, endpoint string) *Story {
	fmt.Println("Querying for story with Endpoint =", endpoint)
	model := &Story{Endpoint: endpoint}
	inserted, err := db.Model(model).Where("endpoint = ?", endpoint).SelectOrInsert(model)
	// db.ModelContext(context.Context(), &models).Select()
	// db.Model(&models).SelectOrInsert()
	// res, err := db.Query(&models, "endpoint = ?", endpoint)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Query Complete. Inserted? %v  Model: %s\n", inserted, model)

	// if res.RowsReturned() == 0 {
	// 	model = &Story{Endpoint: endpoint}
	// 	err := db.Insert(model)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// fmt.Printf("Finished Model: %s\n", model)

	// TODO: foreach story model, do story arcs & pages
	return model
}

func ensureStoryArc(db *pg.DB, story *Story, endpoint string) *StoryArc {
	fmt.Println("Querying for story-arc with Endpoint =", endpoint)
	model := &StoryArc{Endpoint: endpoint, Page: 1}
	err := db.Model(model).Where("endpoint = ?", endpoint).Select()
	if err != nil && err.Error() != "pg: no rows in result set" {
		panic(err)
	}

	fmt.Printf("Query Complete. Model: %s\n", model)

	return model
}

func updateStoryArc(db *pg.DB, arc *StoryArc, page int) *StoryArc {
	fmt.Println("Querying for story-arc with Endpoint =", endpoint)
	model := &StoryArc{Endpoint: endpoint, Page: 1}
	err := db.Model(model).Where("endpoint = ?", endpoint).Select()
	if err != nil && err.Error() != "pg: no rows in result set" {
		panic(err)
	}

	fmt.Printf("Query Complete. Model: %s\n", model)

	return model
}

func main() {
	defer fmt.Println("\n[WORK COMPLETE]")

	// fmt.Printf("Hello, 世界 --- %s\n", time.Now())

	// start := time.Now()
	// response, err := http.Head("https://homestuck.com/story/8131")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Status Code:", response.StatusCode, "Duration:", time.Since(start))

	db := prepareDatabase()
	defer db.Close()

	// stories := lookupStories()
	// fmt.Println("[STORIES]   ", stories)
	// for _, story := range stories {
	// 	ensureStory(db, story)
	story := ensureStory(db, "/bleep")

	// 	storyArcs := lookupStoryArcs(story)
	// 	fmt.Println("[STORY ARCS]", storyArcs)
	// 	for _, arc := range storyArcs {
	// 		arc := ensureStoryArc(db, story, arc)
	arc := ensureStoryArc(db, story, "/bleep/bloop")
	// 		// fmt.Println()
	// 		// fmt.Println("[SEEKING PAGES]")
	// 		latestPage := lookupLatestPage(arc.Endpoint, arc.Page)
	// 		// latestPage := lookupLatestPage("/epilogues/candy", 41)
	// 		// fmt.Printf("\nFound latest page: #%v\n", latestPage)
	updateStoryArc(db, arc, 3)
	// 	}
	// }

}

func uniq(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
