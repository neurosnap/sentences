package sentences

type PunctStrings interface {
	NonPunct() string
	Punctuation() string
	HasSentencePunct(string) bool
}

// Language holds language specific regular expressions to help determine
// information about the text that is being parsed.
type Language struct{}

// Creates a default set of properties for the Language struct
func NewLanguage() *Language {
	return &Language{}
}

func (p *Language) NonPunct() string {
	return `[^\W\d]`
}

func (p *Language) Punctuation() string {
	return ";:,.!?"
}

func (p *Language) HasSentencePunct(text string) bool {
	endPunct := `.!?`
	for _, char := range endPunct {
		for _, achar := range text {
			if char == achar {
				return true
			}
		}
	}

	return false
}
