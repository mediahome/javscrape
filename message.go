package scrape

import "time"

// Message ...
type Message struct {
	Title         string
	OriginalTitle string
	Year          string
	ReleaseDate   time.Time
	ID            string
	Studio        string
	MovieSet      string
	Plot          string
	Genres        []string
	Tags          []string
	Actors        []string
}
