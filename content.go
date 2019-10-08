package scrape

import "time"

// Genre ...
type Genre struct {
	URL     string
	Content string
}

// Content ...
type Content struct {
	From          string //where this
	Uncensored    string
	ID            string
	Title         string
	OriginalTitle string
	Year          string
	ReleaseDate   time.Time
	Studio        string
	MovieSet      string
	Plot          string
	Genres        []*Genre
	Actors        []*Star
	Image         string
	Thumb         string
	Sample        []*Sample
}
