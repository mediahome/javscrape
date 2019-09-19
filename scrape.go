package scrape

type IScrape interface {
}

type scrapeImpl struct {
	grabs []IGrab
}

func NewScrape(grabs ...IGrab) IScrape {
	return scrapeImpl{grabs: grabs}
}
