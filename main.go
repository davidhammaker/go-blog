package main

import (
	"log"
	"net/http"
)

// main maps paths to handlers and starts the server.
func main() {
	http.Handle("/static/index.css", CssHandler{})
	http.Handle("/static/circus.min.css", HighlightCssHandler{})
	http.Handle("/all-posts", EntriesHandler{})
	http.Handle("/file/", FileHandler{})
	http.Handle("/", BlogHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
