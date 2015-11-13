package sentences

import (
	"fmt"
	"regexp"
)

// Groups two adjacent tokens together.
type TokenGrouper interface {
	Group([]*Token) [][2]*Token
}

type DefaultTokenGrouper struct{}

func (p *DefaultTokenGrouper) Group(tokens []*Token) [][2]*Token {
	if len(tokens) == 0 {
		return nil
	}

	pairTokens := make([][2]*Token, 0, len(tokens))

	prevToken := tokens[0]
	for _, tok := range tokens {
		if prevToken == tok {
			continue
		}
		pairTokens = append(pairTokens, [2]*Token{prevToken, tok})
		prevToken = tok
	}
	pairTokens = append(pairTokens, [2]*Token{prevToken, nil})

	return pairTokens
}

// Stores a token of text with annotations produced during sentence boundary detection.
type Token struct {
	Tok         string
	Position    int
	SentBreak   bool
	ParaStart   bool
	LineStart   bool
	Abbr        bool
	periodFinal bool
	reEllipsis  *regexp.Regexp
	reNumeric   *regexp.Regexp
	reInitial   *regexp.Regexp
	reAlpha     *regexp.Regexp
}

func NewToken(token string) *Token {
	tok := Token{
		Tok:        token,
		reEllipsis: regexp.MustCompile(`\.\.+$`),
		reNumeric:  regexp.MustCompile(`-?[\.,]?\d[\d,\.-]*\.?$`),
		reInitial:  regexp.MustCompile(`^[A-Za-z]\.$`),
		reAlpha:    regexp.MustCompile(`^[A-Za-z]+$`),
	}

	return &tok
}

func (p *Token) String() string {
	return fmt.Sprintf("<Token Tok: %q, SentBreak: %t, Abbr: %t, Position: %d>", p.Tok, p.SentBreak, p.Abbr, p.Position)
}
