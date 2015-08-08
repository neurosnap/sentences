package main

import (
	"regexp"
	"strings"
)

// Stores data used to perform sentence boundary detection with punkt
type PunktParameters struct {
	abbrevTypes  map[string]string
	collocations map[string]string
	sentStarters map[string]string
	orthoContext map[string]int
}

func (p *PunktParameters) addOrthoContext(typ string, flag int) {
	p.orthoContext[typ] |= flag
}

// Stores a token of text with annotations produced during sentence boundary detection.
type PunktToken struct {
	reEllipsis  *regexp.Regexp
	reNumeric   *regexp.Regexp
	reInitial   *regexp.Regexp
	reAlpha     *regexp.Regexp
	tok         string
	typ         string
	periodFinal bool
	sentBreak   bool
}

// Returns a case-normalized representation of the token.
func (p *PunktToken) getType(tok string) string {
	return p.reNumeric.ReplaceAllString(tok, "##number##")
}

// The type with its final period removed if it has one.
func (p *PunktToken) typeNoPeriod() string {
	if len(p.typ) > 1 && string(p.typ[len(p.typ)-1]) == "." {
		return string(p.typ[:len(p.typ)-1])
	}
	return p.typ
}

// The type with its final period removed if it is marked as a sentence break.
func (p *PunktToken) typeNoSentPeriod() string {
	if p.sentBreak {
		return p.typeNoPeriod()
	}
	return p.typ
}

// True if the token's first character is uppercase.
func (p *PunktToken) firstUpper() bool {
	firstType := string(p.typ[0])
	return strings.ToUpper(firstType) == firstType
}

// True if the token's first character is lowercase
func (p *PunktToken) firstLower() bool {
	firstType := string(p.typ[0])
	return strings.ToLower(firstType) == firstType
}

// True if the token text is that of an ellipsis.
func (p *PunktToken) firstCase() string {
	if p.firstLower() {
		return "lower"
	} else if p.firstUpper() {
		return "upper"
	}
	return "none"
}

// True if the token text is that of an ellipsis.
func (p *PunktToken) isEllipsis() bool {
	return p.reAlpha.MatchString(p.tok)
}

// True if the token text is that of a number.
func (p *PunktToken) isNumber() bool {
	return string(p.typ[:9]) == "##number##"
}

// True if the token text is that of an initial.
func (p *PunktToken) isInitial() bool {
	return p.reInitial.MatchString(p.tok)
}

// True if the token text is all alphabetic.
func (p *PunktToken) isAlpha() bool {
	return p.reAlpha.MatchString(p.tok)
}

// True if the token is either a number or is alphabetic.
func (p *PunktToken) isNonPunct() bool {
	return re_non_punct.MatchString(p.typ)
}

func NewPunktToken() *PunktToken {
	return &PunktToken{
		reEllipsis: regexp.MustCompile(`\.\.+$`),
		reNumeric:  regexp.MustCompile(`^-?[\.,]?\d[\d,\.-]*\.?$`),
		reInitial:  regexp.MustCompile(`[^\W\d]\.$`),
		reAlpha:    regexp.MustCompile(`[^\W\d]+$`),
	}
}

// Includes common components of PunkTrainer and PunktSentenceTokenizer
type PunktBase struct {
}

// Learns parameters used in Punkt sentence boundary detection
type PunktTrainer struct{}

// A sentence tokenizer which uses an unsupervised algorithm to build a model
// for abbreviation words, collocations, and words that start sentences
// and then uses that model to find sentence boundaries.
type PunktSentenceTokenizer struct{}
