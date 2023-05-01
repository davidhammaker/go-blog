package main

// Entry stores information from a row in the 'entries' table in the
// database.
type Entry struct {
	id          int
	refHost     string
	refPath     string
	description string
}

// blogData stores data to be rendered in the html template for blog
// pages.
type blogData struct {
	BlogTitle   string
	Content     string
	Description string
	Footer      string
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
