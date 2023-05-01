package main

import (
	"html/template"
	"net/http"
)

// CssHandler is an http Handler that serves a static CSS file for the
// blog.
type CssHandler struct{}

// ServeHTTP writes the static CSS file for the CssHandler.
func (CssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.css"))
	w.Header().Set("content-type", "text/css; charset=utf-8")
	tmpl.Execute(w, nil)
}
