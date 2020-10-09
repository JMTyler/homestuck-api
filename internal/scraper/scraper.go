package scraper

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/utils"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

const baseURL = "https://www.homestuck.com"

func FetchAllStories() []map[string]string {
	stories := make([]map[string]string, 0)
	for _, collection := range scrapeCollections() {
		for _, story := range scrapeStories(collection["endpoint"]) {
			story["domain"] = "homestuck.com"
			story["collection"] = collection["title"]
			stories = append(stories, story)
		}
	}
	return stories
}

func SeekLatestPage(_ string, endpoint string, page int) int {
	response, err := http.Head(fmt.Sprintf("%s/%s/%v", baseURL, endpoint, page))
	fmt.Printf("Seeking on /%s -- Checking page %v -- Status Code %v\n", endpoint, page, response.StatusCode)
	if err != nil {
		panic(err)
	}

	if response.StatusCode == 404 {
		return page - 1
	}

	if response.StatusCode == 200 {
		return SeekLatestPage(baseURL, endpoint, page+1)
	}

	panic(fmt.Sprintf("Request to /%s/%v returned unexpected Status Code %v\n", endpoint, page, response.StatusCode))
}

func scrapeCollections() []map[string]string {
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
	return utils.Reverse(result)
}

func scrapeStories(endpoint string) []map[string]string {
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

	links = utils.Uniq(links)

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

	return utils.Reverse(result)
}

func fetch(endpoint string) *goquery.Document {
	response, err := http.Get(baseURL + endpoint)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		panic(err)
	}

	return doc
}
