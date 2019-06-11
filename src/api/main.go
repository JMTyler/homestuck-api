package main

import (
	"encoding/json"
	"fmt"
	"homestuck-watcher/db"
	"homestuck-watcher/fcm"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	fmt.Println()
	defer fmt.Println("\n[[[WORK COMPLETE]]]")
	defer db.CloseDatabase()

	// start := time.Now()
	// time.Since(start)

	http.HandleFunc("/v1/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}

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

	http.HandleFunc("/v1/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}

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
	})

	http.HandleFunc("/v1/stories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}

		storyArcs := new(db.StoryArc).FindAll()
		scrubbed := make([]map[string]interface{}, len(storyArcs))
		for i, arc := range storyArcs {
			scrubbed[i] = arc.Scrub()
		}
		res, _ := json.Marshal(scrubbed)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, string(res))
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.WriteHeader(404)
		fmt.Fprintf(w, "{\"message\":\"Bleep Bloop\"}")
	})

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}
	http.ListenAndServe(":"+port, nil)
}
