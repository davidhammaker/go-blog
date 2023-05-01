package main

import (
	"os"
)

// GetBlogTitle returns the value of the BLOGTITLE environment
// variable, or returns the value "My GO Blog".
func GetBlogTitle() string {
	title, exists := os.LookupEnv("BLOGTITLE")
	if !exists {
		return "My GO Blog"
	}
	return title
}
