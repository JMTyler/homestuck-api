package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	// "io/ioutil"
	"net/http"
	// "time"
	"homestuck-api/db"
	"homestuck-api/fcm"
	"regexp"
	// "sort"
	"encoding/json"
	"io/ioutil"
	"os"
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
		href = regexp.MustCompile("^/log").ReplaceAllString(href, "")

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
		return regexp.MustCompile("/\\d+$").ReplaceAllString(href, "")
	})

	links = uniq(links)

	result := make([]map[string]string, len(links))
	for i, link := range links {
		var title string
		matches := regexp.MustCompile("^/[a-z-]+/([a-z-]+)").FindStringSubmatch(link)
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
	for _, data := range stories {
		fmt.Println("Querying for story with Endpoint =", data["endpoint"])
		story := &db.Story{Endpoint: data["endpoint"], Title: data["title"]}
		story.FindOrCreate()

		storyArcs := lookupStoryArcs(story.Endpoint)
		fmt.Println("[STORY ARCS]", storyArcs)
		for _, data := range storyArcs {
			fmt.Println("Querying for story-arc with Endpoint =", data["endpoint"])
			arc := &db.StoryArc{StoryID: story.ID, Endpoint: data["endpoint"], Title: data["title"], Page: 1}
			arc.FindOrCreate()

			fmt.Println()
			fmt.Println("[SEEKING PAGES]")
			latestPage := lookupLatestPage(arc.Endpoint, arc.Page)
			fmt.Printf("\nFound latest page: #%v\n", latestPage)
			if latestPage != arc.Page {
				arc.ProcessPotato(latestPage)
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
			arc.ProcessPotato(latestPage)
		}
		fmt.Println()
		fmt.Println("----------------------------------------")
		fmt.Println()
	}
}

func populateEmptyStories() {
	stories := lookupStories()
	fmt.Println("[STORIES]", stories)
	for _, data := range stories {
		fmt.Println("Querying for story with Endpoint =", data["endpoint"])
		story := &db.Story{Endpoint: data["endpoint"], Title: data["title"]}
		story.FindOrCreate()

		storyArcs := lookupStoryArcs(story.Endpoint)
		fmt.Println("[STORY ARCS]", storyArcs)
		for _, data := range storyArcs {
			fmt.Println("Querying for story-arc with Endpoint =", data["endpoint"])
			arc := &db.StoryArc{StoryID: story.ID, Endpoint: data["endpoint"], Title: data["title"], Page: 1}
			arc.FindOrCreate()

			fmt.Println()
			fmt.Println("----------------------------------------")
			fmt.Println()
		}
	}
}

func main() {
	fmt.Println()
	defer fmt.Println("\n[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	// slice := lookupStoryArcs("/epilogues")
	// for _, data := range slice {
	// 	fmt.Println("Arc:", data)
	// }

	// new(db.Story).Init()
	// new(db.StoryArc).Init()

	// start := time.Now()
	// time.Since(start)

	// runHeavyweightPoll()
	// runLightweightPoll()

	if len(os.Args) == 1 {
		fmt.Println("No command provided")
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "populate":
		populateEmptyStories()
		return
	case "ping":
		endpoint := "/epilogues/candy"
		if len(os.Args) >= 3 {
			endpoint = os.Args[2]
		}

		arc := &db.StoryArc{Endpoint: endpoint}
		arc.Find()
		fcm.Ping(fcm.SyncEvent, arc.Story.Title, arc.Title, arc.Endpoint, arc.Page)
		return
	case "http":
		// TODO: Might as well prefix with /v1 just in case we ever want it.  Can't hurt!
		http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")

			reqBytes, _ := ioutil.ReadAll(r.Body)
			var req map[string]interface{}
			_ = json.Unmarshal(reqBytes, &req)
			token := req["token"].(string)
			// TODO: Test if this could end up too slow for the web process (once it's being pounded by 1000s of browsers).
			err := fcm.Subscribe([]string{token})
			if err != nil {
				// TODO: Gotta start using log.Fatal() and its ilk.
				fmt.Println(err)
				w.WriteHeader(500)
				fmt.Fprintf(w, "")
				return
			}

			res, _ := json.Marshal(map[string]interface{}{"token": token})
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, string(res))
		})
		http.HandleFunc("/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")

			reqBytes, _ := ioutil.ReadAll(r.Body)
			var req map[string]interface{}
			_ = json.Unmarshal(reqBytes, &req)
			token := req["token"].(string)
			err := fcm.Unsubscribe([]string{token})
			if err != nil {
				// TODO: Gotta start using log.Fatal() and its ilk.
				fmt.Println(err)
				w.WriteHeader(500)
				fmt.Fprintf(w, "")
				return
			}

			fmt.Fprintf(w, "")
		})
		http.HandleFunc("/stories", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")

			// TODO: Should make sure this returns stories in order of creation.
			storyArcs := new(db.StoryArc).FindAll()
			scrubbed := make([]map[string]interface{}, len(storyArcs))
			for i, arc := range storyArcs {
				scrubbed[i] = arc.Scrub()
			}
			res, _ := json.Marshal(scrubbed)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, string(res))
		})
		http.ListenAndServe(":80", nil)
	}

	fmt.Println("Invalid command provided")
}
