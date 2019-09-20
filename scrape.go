package scrape

var debug = false

// IScrape ...
type IScrape interface {
}

type scrapeImpl struct {
	grabs []IGrab
}

// DebugOn ...
func DebugOn() {
	debug = true
}

// NewScrape ...
func NewScrape(grabs ...IGrab) IScrape {
	return scrapeImpl{grabs: grabs}
}

// Find ...
func (scrapeImpl) Find(name string) {

}
