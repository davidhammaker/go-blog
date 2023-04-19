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

type DBHandler struct {
	db *sql.DB
}

func (d *DBHandler) ConnectDB() error {
	addr := os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT")
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   addr,
		DBName: os.Getenv("DBNAME"),
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

type Entry struct {
	id  int
	ref string
}

type CssHandler struct{}

func (CssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.css"))
	w.Header().Set("content-type", "text/css; charset=utf-8")
	tmpl.Execute(w, nil)
}

type blogData struct {
	BlogTitle string
	Content   string
}

func GetBlogTitle() string {
	title, exists := os.LookupEnv("BLOGTITLE")
	if !exists {
		return "My GO Blog"
	}
	return title
}

type BlogHandler struct{}

func (BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[1:]

	var ref string

	if id == "" {
		ref = os.Getenv("HOMEREF")
	} else {

		d := DBHandler{}
		connectErr := d.ConnectDB()
		if connectErr != nil {
			fmt.Println(connectErr)
		}
		row := d.db.QueryRow("SELECT id, ref FROM entries WHERE id = ?", id)
		var ent Entry
		err := row.Scan(&ent.id, &ent.ref)
		if err != nil {
			fmt.Println(err)
		}
		ref = ent.ref
	}

	tmpl := template.Must(template.ParseFiles("static/index.html"))
	w.Header().Set("content-type", "text/html; charset=UTF-8")

	res, resErr := http.Get(ref)
	if resErr != nil {
		fmt.Println(resErr)
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, blogData{BlogTitle: GetBlogTitle(), Content: "# 404\n\nWhatever you are looking for, it's not here."})
		return
	}

	defer res.Body.Close()
	content, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	tmpl.Execute(w, blogData{BlogTitle: GetBlogTitle(), Content: string(content[:])})
}

func main() {
	http.Handle("/static/index.css", CssHandler{})
	http.Handle("/", BlogHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
