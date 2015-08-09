package punkt

import (
	"regexp"
	"strings"
)

// Includes common components of PunkTrainer and PunktSentenceTokenizer
type PunktBase struct {
	// The collection of parameters that determines the behavior of the punkt tokenizer.
	*PunktParameters
	*PunktLanguageVars
}

func NewPunktBase() *PunktBase {
	return &PunktBase{
		PunktLanguageVars: NewPunktLanguageVars(),
	}
}

func (p *PunktBase) AddToken(tokens []*PunktToken, lineTok *WordToken, parastart bool, linestart bool) []*PunktToken {
	nonword := regexp.MustCompile(strings.Join([]string{p.reNonWordChars, p.reMultiCharPunct}, "|"))
	tok := strings.Join([]string{lineTok.First, lineTok.Second}, "")
	if nonword.MatchString(lineTok.Second) || strings.HasSuffix(lineTok.Second, ",") {
		tokOne := &PunktToken{
			Tok:       lineTok.First,
			ParaStart: parastart,
			LineStart: linestart,
		}

		tokTwo := &PunktToken{
			Tok: lineTok.Second,
		}

		tokens = append(tokens, tokOne, tokTwo)
	} else {
		token := &PunktToken{
			Tok:       tok,
			ParaStart: parastart,
			LineStart: linestart,
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func (p *PunktBase) TokenizeWords(text string) []*PunktToken {
	lines := strings.Split(text, "\n")
	tokens := make([]*PunktToken, 0, len(lines))
	parastart := false

	for _, line := range lines {
		if strings.Trim(line, " ") == "" {
			parastart = true
		} else {
			lineToks := p.WordTokenizer(line)

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
func (p *PunktBase) annotateFirstPass(tokens []*PunktToken) []*PunktToken {
	//resultTokens := make([]*PunktToken, 0, len(tokens))
	for _, augTok := range tokens {
		p.firstPassAnnotation(augTok)
	}
	return tokens
}

func (p *PunktBase) firstPassAnnotation(token *PunktToken) {
	tokInEndChars := strings.Index(string(p.sentEndChars), token.Tok)

	if tokInEndChars != -1 {
		token.SentBreak = true
	} else if token.IsEllipsis() {
		token.Ellipsis = true
	} else if token.PeriodFinal && token.Tok[len(token.Tok)-2:] != ".." {
		tokNoPeriod := strings.ToLower(token.Tok[:len(token.Tok)-1])
		tokNoPeriodHypen := strings.Split(tokNoPeriod, "-")
		tokLastHyphEl := string(tokNoPeriodHypen[len(tokNoPeriodHypen)-1])

		if p.PunktParameters.IsAbbr(tokNoPeriod, tokLastHyphEl) {
			token.Abbr = true
		} else {
			token.SentBreak = true
		}
	}
}
