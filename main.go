package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

// DBHandler is a struct whose 'db' attribute is a database Handler.
// The handler can be set with the ConnectDB function.
type DBHandler struct {
	db *sql.DB
}

// ConnectDB connects to a MySQL database, given the following
// environment variables are set: DBHOST, DBPORT, DBUSER, DBPASS,
// DBNAME
func (d *DBHandler) ConnectDB() error {
	addr := os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT")
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	pingErr := db.Ping()
	if pingErr != nil {
		return err
	}
	d.db = db
	return nil
}

// Entry stores information from a row in the 'entries' table in the
// database.
type Entry struct {
	id  int
	ref string
}

// CssHandler is an http Handler that serves a static CSS file for the
// blog.
type CssHandler struct{}

// ServeHTTP writes the static CSS file for the CssHandler.
func (CssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.css"))
	w.Header().Set("content-type", "text/css; charset=utf-8")
	tmpl.Execute(w, nil)
}

// blogData stores data to be rendered in the html template for blog
// pages.
type blogData struct {
	BlogTitle string
	Content   string
}

// GetBlogTitle returns the value of the BLOGTITLE environment
// variable, or returns the value "My GO Blog".
func GetBlogTitle() string {
	title, exists := os.LookupEnv("BLOGTITLE")
	if !exists {
		return "My GO Blog"
	}
	return title
}

// entryData stores data to be rendered for each entry in the list of
// all entries.
type entryData struct {
	Id      int
	Title   string
	Created string
}

// allEntriesData stores all data to be rendered in the list of
// entries.
type allEntriesData struct {
	BlogTitle string
	Entries   []entryData
}

// EntriesHandler is an http Handler that serves an html page containing
// links to all blog entries.
type EntriesHandler struct{}

// EntriesHandler is an http Handler that serves a list of all blog
// entries with links to their respective pages.
func (EntriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := DBHandler{}
	connectErr := d.ConnectDB()
	if connectErr != nil {
		fmt.Println("Could not connect:", connectErr)
	}
	rows, queryErr := d.db.Query("SELECT id, title, created FROM entries ORDER BY created DESC;")
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

	tmpl := template.Must(template.ParseFiles("static/all_posts.html"))
	w.Header().Set("content-type", "text/html; charset=UTF-8")

	tmpl.Execute(w, allEntriesData{BlogTitle: GetBlogTitle(), Entries: entries})
}

// BlogHandler is an http Handler that serves the main html page for
// the blog.
type BlogHandler struct{}

// ServeHTTP fills the html template with data, including the home page
// and any blog entries as derived from database entries.
func (BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[1:]

	var ref string

	if id == "" {
		// If no id is provided, assume this is the home page and use that
		// ref.
		ref = os.Getenv("HOMEREF")

	} else {
		// If an id is provided, attempt to look up the markdown file's URL
		// ref.

		d := DBHandler{}
		connectErr := d.ConnectDB()
		if connectErr != nil {
			fmt.Println("Could not connect:", connectErr)
		}
		row := d.db.QueryRow("SELECT id, ref FROM entries WHERE id = ?", id)
		var ent Entry
		err := row.Scan(&ent.id, &ent.ref)
		if err != nil {
			fmt.Println("Scan failed:", err)
		}
		ref = ent.ref
	}

	tmpl := template.Must(template.ParseFiles("static/index.html"))
	w.Header().Set("content-type", "text/html; charset=UTF-8")

	res, resErr := http.Get(ref)

	// If an error occurs, it means the id was invalid, making this
	// request a 404.
	if resErr != nil {
		fmt.Println("Response error:", resErr)
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, blogData{BlogTitle: GetBlogTitle(), Content: "# 404\n\nWhatever you are looking for, it's not here."})
		return
	}

	defer res.Body.Close()
	content, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Read error:", readErr)
	}

	tmpl.Execute(w, blogData{BlogTitle: GetBlogTitle(), Content: string(content[:])})
}

// main maps paths to handlers and starts the server.
func main() {
	http.Handle("/static/index.css", CssHandler{})
	http.Handle("/all-posts", EntriesHandler{})
	http.Handle("/", BlogHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
