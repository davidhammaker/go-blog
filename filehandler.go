package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// FileHandler is an http Handler that serves files like images and
// videos.
type FileHandler struct{}

// ServeHTTP fetches corresponding files based on the FILEHOST
// environment variable and writes them to the response. Any text-like
// files (HTML, for example) results in a 400 error.
func (FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fileHost := os.Getenv("FILEHOST")
	path := r.URL.Path[6:]

	ref := fileHost + path

	res, resErr := http.Get(ref)
	if resErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer res.Body.Close()

	contentType := res.Header["Content-Type"][0]
	if strings.HasPrefix(contentType, "text") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", res.Header["Content-Type"][0])
	w.Header().Set("content-length", fmt.Sprint(res.ContentLength))

	resBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(resBytes)
	if writeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
