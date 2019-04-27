package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	// "io/ioutil"
	"net/http"
	// "time"
	"./db"
	"regexp"
)

const BaseURL = "https://www.homestuck.com"

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

	if response.StatusCode == 404 {
		return page - 1
	}

	if response.StatusCode == 200 {
		return lookupLatestPage(endpoint, page+1)
	}

	panic(fmt.Sprintf("Request to %s/%v returned unexpected Status Code %v\n", endpoint, page, response.StatusCode))
}

func runHeavyweightPoll() {
	stories := lookupStories()
	fmt.Println("[STORIES]", stories)
	for _, endpoint := range stories {
		fmt.Println("Querying for story with Endpoint =", endpoint)
		story := &db.Story{Endpoint: endpoint}
		story.FindOrCreate()

		storyArcs := lookupStoryArcs(story.Endpoint)
		fmt.Println("[STORY ARCS]", storyArcs)
		for _, endpoint := range storyArcs {
			fmt.Println("Querying for story-arc with Endpoint =", endpoint)
			arc := &db.StoryArc{StoryID: story.ID, Endpoint: endpoint, Page: 1}
			arc.FindOrCreate()

			fmt.Println()
			fmt.Println("[SEEKING PAGES]")
			latestPage := lookupLatestPage(arc.Endpoint, arc.Page)
			fmt.Printf("\nFound latest page: #%v\n", latestPage)
			if latestPage != arc.Page {
				fmt.Printf("Updating story-arc #%v with Page = %v\n", arc.ID, latestPage)
				arc.Page = latestPage
				arc.Update()
			}
			fmt.Println()
			fmt.Println("----------------------------------------")
			fmt.Println()
		}
	}
}

func runLightweightPoll() {
	fmt.Printf("Querying for all story-arcs\n")
	storyArcs := new(db.StoryArc).FindAll()

	fmt.Println("[STORY ARCS]", storyArcs)
	for _, arc := range storyArcs {
		fmt.Println()
		fmt.Println("[SEEKING PAGES]")
		latestPage := lookupLatestPage(arc.Endpoint, arc.Page)
		fmt.Printf("\nFound latest page: #%v\n", latestPage)
		if latestPage != arc.Page {
			fmt.Printf("Updating story-arc #%v with Page = %v\n", arc.ID, latestPage)
			arc.Page = latestPage
			arc.Update()
		}
		fmt.Println()
		fmt.Println("----------------------------------------")
		fmt.Println()
	}
}

func main() {
	defer fmt.Println("[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	// start := time.Now()
	// time.Since(start)

	// runHeavyweightPoll()
	runLightweightPoll()
}
