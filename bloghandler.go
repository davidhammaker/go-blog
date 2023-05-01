package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

// BlogHandler is an http Handler that serves the main html page for
// the blog.
type BlogHandler struct{}

// ServeHTTP fills the html template with data, including the home page
// and any blog entries as derived from database entries.
func (BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[1:]

	var ref string
	var description string

	if id == "" {
		// If no id is provided, assume this is the home page and use that
		// ref.
		ref = os.Getenv("HOMEREF")
		description = os.Getenv("HOMEDESCRIPTION")

	} else {
		// If an id is provided, attempt to look up the markdown file's URL
		// ref.

		db, connectErr := ConnectDB()
		if connectErr != nil {
			fmt.Println("Could not connect:", connectErr)
		}
		row := db.QueryRow("SELECT id, refHost, refPath, description FROM entries WHERE id = ?", id)
		var ent Entry
		err := row.Scan(&ent.id, &ent.refHost, &ent.refPath, &ent.description)
		if err != nil {
			fmt.Println("Scan failed:", err)
		}
		ref = ent.refHost + ent.refPath
		description = ent.description
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	w.Header().Set("content-type", "text/html; charset=UTF-8")

	res, resErr := http.Get(ref)

	// If an error occurs, it means the id was invalid, making this
	// request a 404.
	if resErr != nil {
		fmt.Println("Response error:", resErr)
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, blogData{BlogTitle: GetBlogTitle(), Content: "# 404\n\nWhatever you are looking for, it's not here.", Description: "404 NOT FOUND", Footer: ""})
		return
	}

	defer res.Body.Close()
	content, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Read error:", readErr)
	}

	tmpl.Execute(w, blogData{BlogTitle: GetBlogTitle(), Content: string(content[:]), Description: description, Footer: os.Getenv("FOOTER")})
}
