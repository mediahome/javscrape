package scrape

import "time"

// Message ...
type Message struct {
	ID            string
	Title         string
	OriginalTitle string
	Year          string
	ReleaseDate   time.Time
	Studio        string
	MovieSet      string
	Plot          string
	Genres        []string
	Actors        []string
}

// Idols ...
type Idols struct {
	Image    string
	StarLink string
	Name     string
}
