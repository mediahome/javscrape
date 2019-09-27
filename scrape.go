package scrape

var debug = false

// IScrape ...
type IScrape interface {
	GrabSample(bool)
}

type scrapeImpl struct {
	grabs  []IGrab
	sample bool
}

// GrabSample ...
func (impl *scrapeImpl) GrabSample(b bool) {
	impl.sample = b
}

// DebugOn ...
func DebugOn() {
	debug = true
}

// NewScrape ...
func NewScrape(grabs ...IGrab) IScrape {
	return &scrapeImpl{grabs: grabs}
}

// Find ...
func (impl *scrapeImpl) Find(name string) (msg []*Message, e error) {
	msg = *new([]*Message)
	for _, grab := range impl.grabs {
		iGrab, e := grab.Find(name)
		if e != nil {
			log.With("name", grab.Name()).Error(e)
			continue
		}
		e = iGrab.Decode(msg)
		if e != nil {
			log.With("name", grab.Name()).Error(e)
		}
	}
	return
}
