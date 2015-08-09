package punkt

import (
	"regexp"
	"strings"
)

// Stores a token of text with annotations produced during sentence boundary detection.
type PunktToken struct {
	reEllipsis  *regexp.Regexp
	reNumeric   *regexp.Regexp
	reInitial   *regexp.Regexp
	reAlpha     *regexp.Regexp
	Tok         string
	Typ         string
	PeriodFinal bool
	SentBreak   bool
	ParaStart   bool
	LineStart   bool
	Ellipsis    bool
	Abbr        bool
}

// Returns a case-normalized representation of the token.
func (p *PunktToken) getType(tok string) string {
	return p.reNumeric.ReplaceAllString(tok, "##number##")
}

// The type with its final period removed if it has one.
func (p *PunktToken) TypeNoPeriod() string {
	if len(p.Typ) > 1 && string(p.Typ[len(p.Typ)-1]) == "." {
		return string(p.Typ[:len(p.Typ)-1])
	}
	return p.Typ
}

// The type with its final period removed if it is marked as a sentence break.
func (p *PunktToken) TypeNoSentPeriod() string {
	if p.SentBreak {
		return p.TypeNoPeriod()
	}
	return p.Typ
}

// True if the token's first character is uppercase.
func (p *PunktToken) FirstUpper() bool {
	if p.Tok == "" {
		return false
	}
	firstTok := string(p.Tok[0])
	return strings.ToUpper(firstTok) == firstTok
}

// True if the token's first character is lowercase
func (p *PunktToken) FirstLower() bool {
	if p.Tok == "" {
		return false
	}
	firstTok := string(p.Tok[0])
	return strings.ToLower(firstTok) == firstTok
}

// True if the token text is that of an ellipsis.
func (p *PunktToken) FirstCase() string {
	if p.FirstLower() {
		return "lower"
	} else if p.FirstUpper() {
		return "upper"
	}
	return "none"
}

// True if the token text is that of an ellipsis.
func (p *PunktToken) IsEllipsis() bool {
	return p.reAlpha.MatchString(p.Tok)
}

// True if the token text is that of a number.
func (p *PunktToken) IsNumber() bool {
	return strings.HasPrefix(p.Tok, "##number##")
}

// True if the token text is that of an initial.
func (p *PunktToken) IsInitial() bool {
	return p.reInitial.MatchString(p.Tok)
}

// True if the token text is all alphabetic.
func (p *PunktToken) IsAlpha() bool {
	return p.reAlpha.MatchString(p.Tok)
}

// True if the token is either a number or is alphabetic.
func (p *PunktToken) IsNonPunct() bool {
	return ReNonPunct.MatchString(p.Typ)
}

func NewPunktToken(token string) *PunktToken {
	tok := PunktToken{
		Tok:        token,
		reEllipsis: regexp.MustCompile(`\.\.+$`),
		reNumeric:  regexp.MustCompile(`^-?[\.,]?\d[\d,\.-]*\.?$`),
		reInitial:  regexp.MustCompile(`[^\W\d]\.$`),
		reAlpha:    regexp.MustCompile(`[^\W\d]+$`),
	}

	tok.Typ = tok.getType(token)
	tok.PeriodFinal = strings.HasSuffix(token, ".")

	return &tok
}
