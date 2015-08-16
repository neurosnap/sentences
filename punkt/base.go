package punkt

import (
	"regexp"
	"strings"
)

// Includes common components of Trainer and SentenceTokenizer
type Base struct {
	// The collection of parameters that determines the behavior of the punkt tokenizer.
	*Storage
	*Language
}

func NewBase() *Base {
	return &Base{
		Storage:  NewStorage(),
		Language: NewLanguage(),
	}
}

func (p *Base) AddToken(tokens []*Token, pairTok *PairToken, parastart bool, linestart bool) []*Token {
	nonword := regexp.MustCompile(strings.Join([]string{p.reNonWordChars, ReMultiCharPunct}, "|"))
	tok := strings.Join([]string{pairTok.First, pairTok.Second}, "")

	if nonword.MatchString(pairTok.Second) || strings.HasSuffix(pairTok.Second, ",") {
		tokOne := NewToken(pairTok.First)
		tokOne.ParaStart = parastart
		tokOne.LineStart = linestart

		tokTwo := NewToken(pairTok.Second)

		tokens = append(tokens, tokOne, tokTwo)
	} else {
		token := NewToken(tok)
		token.ParaStart = parastart
		token.LineStart = linestart
		tokens = append(tokens, token)
	}

	return tokens
}

func (p *Base) TokenizeWords(text string) []*Token {
	lines := strings.Split(text, "\n")
	tokens := make([]*Token, 0, len(lines))
	parastart := false

	for _, line := range lines {
		if strings.Trim(line, " ") == "" || line == " " {
			parastart = true
		} else {
			lineToks := WordTokenizer(line)
			for index, lineTok := range lineToks {
				if index == 0 {
					tokens = p.AddToken(tokens, lineTok, parastart, true)
					parastart = false
				} else {
					tokens = p.AddToken(tokens, lineTok, parastart, false)
				}
			}
		}
	}

	return tokens
}

/*
Perform the first pass of annotation, which makes decisions
based purely based on the word type of each word:
	- '?', '!', and '.' are marked as sentence breaks.
	- sequences of two or more periods are marked as ellipsis.
	- any word ending in '.' that's a known abbreviation is marked as an abbreviation.
	- any other word ending in '.' is marked as a sentence break.

Return these annotations as a tuple of three sets:
	- sentbreak_toks: The indices of all sentence breaks.
	- abbrev_toks: The indices of all abbreviations.
	- ellipsis_toks: The indices of all ellipsis marks.
*/
func (p *Base) annotateFirstPass(tokens []*Token) []*Token {
	for _, augTok := range tokens {
		p.firstPassAnnotation(augTok)
	}
	return tokens
}

func (p *Base) firstPassAnnotation(token *Token) {
	tokInEndChars := strings.Index(string(p.sentEndChars), token.Tok)

	if tokInEndChars != -1 {

		token.SentBreak = true
	} else if token.IsEllipsis() {
		token.Ellipsis = true
	} else if token.PeriodFinal && !strings.HasSuffix(token.Tok, "..") {
		tokNoPeriod := strings.ToLower(token.Tok[:len(token.Tok)-1])
		tokNoPeriodHypen := strings.Split(tokNoPeriod, "-")
		tokLastHyphEl := string(tokNoPeriodHypen[len(tokNoPeriodHypen)-1])

		if p.IsAbbr(tokNoPeriod, tokLastHyphEl) {
			token.Abbr = true
		} else {
			token.SentBreak = true
		}
	}
}

func (p *Base) pairIter(tokens []*Token) [][2]*Token {
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
