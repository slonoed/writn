package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	mgo "gopkg.in/mgo.v2"
)

type record struct {
	URL         string `json:"url"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

func main() {
	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	ms, err := mgo.Dial(dbURL)
	if err != nil {
		log.Panicf("Dial mongo error: %s. URL: %s", err, dbURL)
	}
	defer ms.Close()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	http.HandleFunc("/api/addRecord", func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Err", 500)
			return
		}

		// Unmarshal
		var rec record
		err = json.Unmarshal(b, &rec)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Wrong format", 500)
			return
		}

		if rec.URL == "" {
			http.Error(w, "URL should not be empty", 500)
			return
		}

		ms.DB(dbName).C("records").Insert(&rec)

		fmt.Fprintf(w, "ok2")
	})

	http.ListenAndServe(":8080", nil)

	log.Println("Server running!")
	// Shut down when SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}
