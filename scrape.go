package scrape

// IScrape ...
type IScrape interface {
}

type scrapeImpl struct {
	grabs []IGrab
}

// NewScrape ...
func NewScrape(grabs ...IGrab) IScrape {
	return scrapeImpl{grabs: grabs}
}

// Find ...
func (scrapeImpl) Find(name string) {

}
