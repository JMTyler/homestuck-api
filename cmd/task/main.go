package main

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/src/db"
	"github.com/JMTyler/homestuck-watcher/src/fcm"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"regexp"
	"strings"
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

func reverse(slice []map[string]string) []map[string]string {
	result := make([]map[string]string, len(slice))
	for i := 0; i < len(slice); i++ {
		result[i] = slice[len(slice)-i-1]
	}
	return result
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

func lookupStories() []map[string]string {
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
	})

	result := make([]map[string]string, links.Size())
	links.Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = regexp.MustCompile("^/log/").ReplaceAllString(href, "")

		title, _ := s.Parent().Parent().Find("h2").Html()

		entry := make(map[string]string)
		entry["endpoint"] = href
		entry["title"] = strings.Title(strings.ToLower(title))
		result[i] = entry

		fmt.Printf("HTML(STORY):  %s  --  %s\n", entry["title"], entry["endpoint"])
	})

	// TODO: Make result implement sort.Interface so we can use sort.Reverse() here.
	return reverse(result)
}

func lookupStoryArcs(endpoint string) []map[string]string {
	doc := fetch("/log/" + endpoint)

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
		return strings.TrimPrefix(regexp.MustCompile("/\\d+$").ReplaceAllString(href, ""), "/")
	})

	links = uniq(links)

	result := make([]map[string]string, len(links))
	for i, link := range links {
		var title string
		matches := regexp.MustCompile("^[a-z-]+/([a-z-]+)").FindStringSubmatch(link)
		if matches != nil {
			title = strings.Title(strings.ReplaceAll(matches[1], "-", " "))
		}

		entry := make(map[string]string)
		entry["endpoint"] = link
		entry["title"] = title
		result[i] = entry

		fmt.Printf("HTML(ARC):  %v  --  %s\n", entry["title"], entry["endpoint"])
	}

	return reverse(result)
}

func populateEmptyStories() {
	stories := lookupStories()
	for _, data := range stories {
		collection := data["title"]

		storyArcs := lookupStoryArcs(data["endpoint"])
		for _, data := range storyArcs {
			story := &db.Story{
				Domain: "homestuck.com",
				Endpoint: data["endpoint"],
				Collection: collection,
				Title: data["title"],
				Page: 1,
			}
			story.FindOrCreate()
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No command provided")
		return
	}

	fmt.Println()
	defer fmt.Println("\n[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	switch os.Args[1] {
	case "populate":
		populateEmptyStories()
		return
	case "ping":
		endpoint := "epilogues/candy"
		if len(os.Args) >= 3 {
			endpoint = os.Args[2]
		}

		story := &db.Story{Domain: "homestuck.com", Endpoint: endpoint}
		story.Find()
		fcm.Ping(fcm.SyncEvent, story.Collection, story.Title, story.Domain, story.Endpoint, story.Page)
		return
	case "potato":
		endpoint := "epilogues/candy"
		if len(os.Args) >= 3 {
			endpoint = os.Args[2]
		}

		story := &db.Story{Domain: "homestuck.com", Endpoint: endpoint}
		story.Find()
		fcm.Ping(fcm.PotatoEvent, story.Collection, story.Title, story.Domain, story.Endpoint, story.Page)
		return
	}

	fmt.Println("Invalid command provided")
}
