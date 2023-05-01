package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// EntriesHandler is an http Handler that serves an html page containing
// links to all blog entries.
type EntriesHandler struct{}

// EntriesHandler is an http Handler that serves a list of all blog
// entries with links to their respective pages.
func (EntriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db, connectErr := ConnectDB()
	if connectErr != nil {
		fmt.Println("Could not connect:", connectErr)
	}
	rows, queryErr := db.Query("SELECT id, title, created FROM entries ORDER BY created DESC;")
	if queryErr != nil {
		fmt.Println("SELECT failed:", queryErr)
	}
	defer rows.Close()
	var entries []entryData
	for rows.Next() {
		var ent entryData
		err := rows.Scan(&ent.Id, &ent.Title, &ent.Created)
		if err != nil {
			fmt.Println("Scan failed:", err)
		}
		entries = append(entries, ent)
	}

	tmpl := template.Must(template.ParseFiles("templates/all_posts.html"))
	w.Header().Set("content-type", "text/html; charset=UTF-8")

	tmpl.Execute(w, allEntriesData{BlogTitle: GetBlogTitle(), Entries: entries})
}
