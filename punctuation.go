package sentences

type PunctStrings interface {
	NonPunct() string
	Punctuation() string
	HasSentencePunct(string) bool
}

type DefaultPunctStrings struct{}

// Creates a default set of properties
func NewPunctStrings() *DefaultPunctStrings {
	return &DefaultPunctStrings{}
}

func (p *DefaultPunctStrings) NonPunct() string {
	return `[^\W\d]`
}

func (p *DefaultPunctStrings) Punctuation() string {
	return ";:,.!?"
}

func (p *DefaultPunctStrings) HasSentencePunct(text string) bool {
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
