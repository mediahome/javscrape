package scrape

// IGrab ...
type IGrab interface {
	Find(string) error
	Decode(*Message) error
}
