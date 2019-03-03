package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/olivere/elastic"
)

type application struct {
	conn *elastic.Client
}

type result struct {
	Expires   time.Time
	Timestamp time.Time
	Title     string
	Created   time.Time
	ID        int
	Content   string
	Version   string
}

func main() {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Println(err)
		return
	}

	app := &application{conn: client}

	mux := http.NewServeMux()
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/search", app.searchKeyword)

	port := ":8080"
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("homepage.tmpl.html")
	if err != nil {
		log.Println(err)
		return
	}

	t.Execute(w, nil)
}

func (app *application) searchKeyword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	q := elastic.NewQueryStringQuery("*" + r.FormValue("keywords") + "*")
	indices := []string{"poems"}
	results, err := app.queryES(q, indices...)
	if err != nil {
		log.Println(err)
		return
	}

	if len(results) == 0 {
		w.Write([]byte("Sorry. No results were found."))
		return
	}

	for _, v := range results {
		fmt.Fprintf(w, "Title: %s\n", v.Title)
		fmt.Fprintf(w, "%s\n\n\n", v.Content)
	}
}

func (app *application) queryES(query elastic.Query, indices ...string) ([]result, error) {
	ctx := context.Background()

	searchResult, err := app.conn.Search(indices...).
		Query(query).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var results []result

	for _, hit := range searchResult.Hits.Hits {
		r := &result{}
		err := json.Unmarshal(*hit.Source, r)
		if err != nil {
			return nil, err
		}
		results = append(results, *r)
	}

	return results, nil
}
