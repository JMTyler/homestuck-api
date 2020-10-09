package main

import (
	"encoding/json"
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/fcm"
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
		if req["token"] == nil || req["token"] == false || req["token"] == "" {
			w.WriteHeader(422)
			fmt.Fprintf(w, "Required field `token` was empty")
			return
		}

		token := req["token"].(string)
		err := fcm.Subscribe("v1", token)
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
		if req["token"] == nil || req["token"] == false || req["token"] == "" {
			w.WriteHeader(422)
			fmt.Fprintf(w, "Required field `token` was empty")
			return
		}

		token := req["token"].(string)
		err := fcm.Unsubscribe("v1", token)
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

		stories := new(db.Story).FindAll("v1")
		scrubbed := make([]map[string]interface{}, len(stories))
		for i, model := range stories {
			scrubbed[i] = model.Scrub("v1")
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
