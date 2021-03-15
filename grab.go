package scrape

// IGrab ...
type IGrab interface {
	MainPage(url string)
	SetSample(bool)
	SetExact(bool)
	SetLanguage(language GrabLanguage)
	Name() string
	Find(string) (IGrab, error)
	HasNext() bool
	Next() (IGrab, error)
	Result() ([]Content, error)
	SetForce(force bool)
}

// GrabLanguage ...
type GrabLanguage int

// GrabLanguage detail ...
const (
	LanguageEnglish GrabLanguage = iota
	LanguageJapanese
	LanguageChineseSimple
	LanguageChineseTraditional
	LanguageKorea
)

var languageGrabStringList = map[GrabLanguage]string{
	LanguageEnglish:            "english",
	LanguageJapanese:           "japanese",
	LanguageChineseSimple:      "simple chinese",
	LanguageChineseTraditional: "traditional chinese",
	LanguageKorea:              "korea",
}

func (g GrabLanguage) String() string {
	v, b := languageGrabStringList[g]
	if !b {
		return ""
	}
	return v
}
