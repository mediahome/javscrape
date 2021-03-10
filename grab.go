package scrape

// IGrab ...
type IGrab interface {
	MainPage(url string)
	SetSample(bool)
	SetExact(bool)
	Name() string
	Find(string) (IGrab, error)
	HasNext() bool
	Next() (IGrab, error)
	Result() ([]*Content, error)
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
