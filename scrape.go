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
func (impl *scrapeImpl) Find(name string) {
	for _, grab := range impl.grabs {
		iGrab, e := grab.Find(name)
		if e != nil {

		}
	}
}
