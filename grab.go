package scrape

// IGrab ...
type IGrab interface {
	Find(string) (IGrab, error)
	Decode(*Message) error
}

// GrabLanguage ...
type GrabLanguage int

// GrabLanguage detail ...
const (
	LanguageEnglish GrabLanguage = iota
	LanguageJapanese
	LanguageChinese
	LanguageKorea
)
