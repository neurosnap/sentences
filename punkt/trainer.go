package punkt

import (
	"fmt"
	"github.com/neurosnap/go-sentences/utils"
	"math"
	"strings"
)

type AbbrevType struct {
	Typ   string
	Score float64
	IsAdd bool
}

/*
The following constants are used to describe the orthographic
contexts in which a word can occur.  BEG=beginning, MID=middle,
UNK=unknown, UC=uppercase, LC=lowercase, NC=no case.
*/
const (
	// Beginning of a sentence with upper case.
	orthoBegUc = 1 << 1
	// Middle of a sentence with upper case.
	orthoMidUc = 1 << 2
	// Unknown position in a sentence with upper case.
	orthoUnkUc = 1 << 3
	// Beginning of a sentence with lower case.
	orthoBegLc = 1 << 4
	// Middle of a sentence with lower case.
	orthoMidLc = 1 << 5
	// Unknown position in a sentence with lower case.
	orthoUnkLc = 1 << 6
	// Occurs with upper case.
	orthoUc = orthoBegUc + orthoMidUc + orthoUnkUc
	// Occurs with lower case.
	orthoLc = orthoBegLc + orthoMidLc + orthoUnkLc
)

/*
A map from context position and first-letter case to the
appropriate orthographic context flag.
*/
var orthoMap = map[[2]string]int{
	[2]string{"initial", "upper"}:  orthoBegUc,
	[2]string{"internal", "upper"}: orthoMidUc,
	[2]string{"unknown", "upper"}:  orthoUnkUc,
	[2]string{"initial", "lower"}:  orthoBegLc,
	[2]string{"internal", "lower"}: orthoMidLc,
	[2]string{"unknown", "lower"}:  orthoUnkLc,
}

func boolToFloat64(cond bool) float64 {
	if cond {
		return 1
	}
	return 0
}

// Learns parameters used in Punkt sentence boundary detection
type PunktTrainer struct {
	*PunktBase
	typeFreqDist         *utils.FreqDist
	collocationFreqDist  *utils.FreqDist
	sentStarterFreqDist  *utils.FreqDist
	sentBreakCount       float64
	numPeriodToks        float64
	Abbrev               float64
	Collocation          float64
	SentStarter          float64
	MinCollocFreq        float64
	IgnoreAbbrevPenalty  bool
	finalized            bool
	IncludeAllCollocs    bool
	IncludeAbbrevCollocs bool
	AbbrevBackoff        int
}

func NewPunktTrainer(trainText string) *PunktTrainer {
	trainer := &PunktTrainer{
		PunktBase:           NewPunktBase(),
		typeFreqDist:        utils.NewFreqDist(),
		collocationFreqDist: utils.NewFreqDist(),
		sentStarterFreqDist: utils.NewFreqDist(),
		finalized:           true,
		Abbrev:              0.3,
		AbbrevBackoff:       5,
		Collocation:         7.88,
		SentStarter:         30,
		MinCollocFreq:       1,
		IncludeAllCollocs:   true,
	}

	if trainText != "" {
		trainer.Train(trainText, true)
	}

	return trainer
}

func (p *PunktTrainer) Train(text string, finalize bool) {
	p.trainTokens(p.TokenizeWords(text))
	if finalize {
		p.FinalizeTraining()
	}
}

func (p *PunktTrainer) trainTokens(tokens []*PunktToken) {
	p.finalized = false
	/*
		Find the frequency of each case-normalized type.  (Don't
		strip off final periods.)  Also keep track of the number of
		tokens that end in periods.
	*/
	for _, tok := range tokens {
		if p.typeFreqDist.Samples[tok.Typ] == 0 {
			p.typeFreqDist.Samples[tok.Typ] = 1
		} else {
			p.typeFreqDist.Samples[tok.Typ] += 1
		}

		if tok.PeriodFinal {
			p.numPeriodToks += 1
		}
	}

	// Look for new abbreviations, and for types that no longer are
	uniqueTypes := p.uniqueTypes(tokens)
	for _, abbrType := range p.reclassifyAbbrevTypes(uniqueTypes) {
		if abbrType.Score >= p.Abbrev {
			if abbrType.IsAdd {
				p.PunktParameters.AbbrevTypes.Add(abbrType.Typ)
			}
		} else {
			if !abbrType.IsAdd {
				p.PunktParameters.AbbrevTypes.Remove(abbrType.Typ)
			}
		}
	}

	/*
		Make a preliminary pass through the document, marking likely
		sentence breaks, abbreviations, and ellipsis tokens.
	*/
	fpTokens := p.annotateFirstPass(tokens)

	// Check what contexts each word type can appear in, given the case of its first letter.
	p.getOrthographData(fpTokens)

	// We need total number of sentence breaks to find sentence starters
	p.sentBreakCount += p.getSentBreakCount(tokens)

	for _, tokPair := range p.pairIter(tokens) {
		if !tokPair[0].PeriodFinal || tokPair[1] == nil {
			continue
		}

		if p.isRareAbbrevType(tokPair[0], tokPair[1]) {
			p.PunktParameters.AbbrevTypes.Add(tokPair[0].TypeNoPeriod())
		}

		if p.isPotentialSentStarter(tokPair[1], tokPair[0]) {
			p.sentStarterFreqDist.Samples[tokPair[1].Typ] += 1
		}

		if p.isPotentialCollocation(tokPair[0], tokPair[1]) {
			var bigramColl = strings.Join([]string{
				tokPair[0].TypeNoPeriod(),
				tokPair[1].TypeNoSentPeriod(),
			}, ",")
			p.collocationFreqDist.Samples[bigramColl] += 1
		}
	}
}

func (p *PunktTrainer) pairIter(tokens []*PunktToken) [][2]*PunktToken {
	pairTokens := make([][2]*PunktToken, 0, len(tokens))

	prevToken := tokens[0]
	for _, tok := range tokens {
		pairTokens = append(pairTokens, [2]*PunktToken{prevToken, tok})
		prevToken = tok
	}
	pairTokens = append(pairTokens, [2]*PunktToken{prevToken, nil})

	return pairTokens
}

func (p *PunktTrainer) uniqueTypes(tokens []*PunktToken) []string {
	unique := NewSetString()

	for _, tok := range tokens {
		unique.Add(tok.Typ)
	}

	return unique.Array()
}

func (p *PunktTrainer) TrainTokens(tokens []string, finalize bool) {}

/*
Uses data that has been gathered in training to determine likely
collocations and sentence starters.
*/
func (p *PunktTrainer) FinalizeTraining() {
	p.PunktParameters.SentStarters.items = map[string]int{}
	for _, ss := range p.findSentStarters() {
		p.PunktParameters.SentStarters.Add(ss.Typ)
	}

	p.PunktParameters.Collocations.items = map[string]int{}
	for _, val := range p.findCollocations() {
		p.PunktParameters.Collocations.Add(strings.Join([]string{val.TypOne, val.TypTwo}, ","))
	}

	p.finalized = true
}

/*
(Re)classifies each given token if
	- it is period-final and not a known abbreviation; or
	- it is not period-final and is otherwise a known abbreviation
by checking whether its previous classification still holds according
to the heuristics of section 3.
Yields triples (abbr, score, is_add) where abbr is the type in question,
score is its log-likelihood with penalties applied, and is_add specifies
whether the present type is a candidate for inclusion or exclusion as an
abbreviation, such that:
	- (is_add and score >= 0.3)    suggests a new abbreviation; and
	- (not is_add and score < 0.3) suggests excluding an abbreviation.
*/
func (p *PunktTrainer) reclassifyAbbrevTypes(types []string) []*AbbrevType {
	abbrTypes := make([]*AbbrevType, 0, len(types))

	for _, typ := range types {
		// Check some basic conditions, to rule out words that are
		// clearly not abbrev_types.
		isPunct := ReNonPunct.FindString(typ) == ""
		if isPunct || typ == "##number##" {
			continue
		}

		var isAdd bool
		if strings.HasSuffix(typ, ".") {
			if p.PunktParameters.AbbrevTypes.Has(typ) {
				continue
			}
			typ = typ[:len(typ)-1]
			isAdd = true
		} else {
			if !p.PunktParameters.AbbrevTypes.Has(typ) {
				continue
			}
			isAdd = false
		}

		numPeriods := float64(strings.Count(typ, ".") + 1)
		numNonPeriods := float64(float64(len(typ)) - numPeriods + 1)
		/*
			Let <a> be the candidate without the period, and <b>
			be the period.  Find a log likelihood ratio that
			indicates whether <ab> occurs as a single unit (high
			value of ll), or as two independent units <a> and
			<b> (low value of ll).
		*/
		typPeriod := strings.Join([]string{typ, "."}, "")
		countWithPeriod := float64(p.typeFreqDist.Samples[typPeriod])
		countWithoutPeriod := float64(p.typeFreqDist.Samples[typ])
		/*
			Apply three scaling factors to 'tweak' the basic log
			likelihood ratio:
				F_length: long word -> less likely to be an abbrev
				F_periods: more periods -> more likely to be an abbrev
				F_penalty: penalize occurrences w/o a period
		*/
		likely := p.dunningLogLikelihood(
			countWithPeriod+countWithoutPeriod,
			p.numPeriodToks, countWithPeriod,
			p.typeFreqDist.N(),
		)
		fLength := math.Exp(-numNonPeriods)
		var fPenalty float64
		if p.IgnoreAbbrevPenalty {
			fPenalty = 1
		} else {
			fPenalty = math.Pow(numNonPeriods, -countWithoutPeriod)
		}
		score := likely * fLength * numPeriods * fPenalty

		//fmt.Println(typ, likely)
		abbrTypes = append(abbrTypes, &AbbrevType{typ, score, isAdd})
	}

	return abbrTypes
}

/*
 A function that calculates the modified Dunning log-likelihood
 ratio scores for abbreviation candidates.  The details of how
 this works is available in the paper.
*/
func (p *PunktTrainer) dunningLogLikelihood(countA, countB, countAB, N float64) float64 {
	p1 := countB / N
	p2 := 0.99

	nullHypo := (countAB*math.Log(p1) + (countA-countAB)*math.Log(1.0-p1))
	altHypo := (countAB*math.Log(p2) + (countA-countAB)*math.Log(1.0-p2))

	likelihood := nullHypo - altHypo

	return -2.0 * likelihood
}

/*
Collect information about whether each token type occurs
with different case patterns (i) overall, (ii) at
sentence-initial positions, and (iii) at sentence-internal
positions.
*/
func (p *PunktTrainer) getOrthographData(tokens []*PunktToken) {
	context := "internal"

	/*
		If we encounter a paragraph break, then it's a good sign
		that it's a sentence break.  But err on the side of
		caution (by not positing a sentence break) if we just
		saw an abbreviation.
	*/
	for _, tok := range tokens {
		if tok.ParaStart && context != "unknown" {
			context = "initial"
		}

		/*
			If we're at the beginning of a line, then we can't decide
			between 'internal' and 'initial'.
		*/
		if tok.LineStart && context == "internal" {
			context = "unknown"
		}
		/*
		   Find the case-normalized type of the token.  If it's a
		   sentence-final token, strip off the period.
		*/
		typ := tok.TypeNoSentPeriod()

		// Update the orthographic context table
		flag := orthoMap[[2]string{context, tok.FirstCase()}]
		if flag != 0 {
			fmt.Println(typ, context, tok.FirstCase())
			p.PunktParameters.addOrthoContext(typ, flag)
		}

		// Decide whether the next word is at a sentence boundary.
		if tok.SentBreak {
			if !(tok.IsNumber() || tok.IsInitial()) {
				context = "initial"
			} else {
				context = "unknown"
			}
		} else if tok.Ellipsis || tok.Abbr {
			//fmt.Println(tok.Ellipsis, tok.Abbr)
			context = "unknown"
		} else {
			context = "internal"
		}
	}
}

/*
Returns the number of sentence breaks marked in a given set of
augmented tokens.
*/
func (p *PunktTrainer) getSentBreakCount(tokens []*PunktToken) float64 {
	sum := 0.0

	for _, tok := range tokens {
		if tok.SentBreak {
			sum += 1
		}
	}

	return sum
}

/*
This function combines the work done by the original code's
functions `count_orthography_context`, `get_orthography_count`,
and `get_rare_abbreviations`.
*/
func (p *PunktTrainer) isRareAbbrevType(curTok, nextTok *PunktToken) bool {
	/*
		A word type is counted as a rare abbreviation if:
			- it's not already marked as an abbreviation
			- it occurs fewer than ABBREV_BACKOFF times
			- either it is followed by a sentence-internal punctuation
			mark, *or* it is followed by a lower-case word that
			sometimes appears with upper case, but never occurs with
			lower case at the beginning of sentences.
	*/
	if curTok.Abbr || !curTok.SentBreak {
		return false
	}

	/*
		Find the case-normalized type of the token.  If it's
		a sentence-final token, strip off the period.
	*/
	typ := curTok.TypeNoSentPeriod()

	/*
		Proceed only if the type hasn't been categorized as an
		abbreviation already, and is sufficiently rare...
	*/
	count := p.typeFreqDist.Samples[typ] + p.typeFreqDist.Samples[typ[:len(typ)-1]]
	if p.PunktParameters.AbbrevTypes.Has(typ) || count >= p.AbbrevBackoff {
		return false
	}

	/*
		Record this token as an abbreviation if the next
		token is a sentence-internal punctuation mark.
		[XX] :1 or check the whole thing??
	*/
	if strings.Contains(p.PunktLanguageVars.internalPunctuation, nextTok.Tok[:1]) {
		return true
	}

	/*
		Record this type as an abbreviation if the next
		token...  (i) starts with a lower case letter,
		(ii) sometimes occurs with an uppercase letter,
		and (iii) never occus with an uppercase letter
		sentence-internally.
		[xx] should the check for (ii) be modified??
	*/
	if nextTok.FirstLower() {
		typTwo := nextTok.TypeNoSentPeriod()
		typTwoOrthoCtx := p.PunktParameters.OrthoContext.items[typTwo]

		if (typTwoOrthoCtx&orthoBegUc) == 1 && (typTwoOrthoCtx&orthoMidUc) == 0 {
			return true
		}
	}

	return false
}

/*
Returns True given a token and the token that preceds it if it
seems clear that the token is beginning a sentence.
*/
func (p *PunktTrainer) isPotentialSentStarter(curTok, prevTok *PunktToken) bool {
	return prevTok.SentBreak && !(prevTok.IsNumber() || prevTok.IsInitial()) && curTok.IsAlpha()
}

/*
Returns True if the pair of tokens may form a collocation given
log-likelihood statistics.
*/
func (p *PunktTrainer) isPotentialCollocation(tokOne, tokTwo *PunktToken) bool {
	return ((p.IncludeAllCollocs || p.IncludeAbbrevCollocs && tokOne.Abbr) ||
		(tokOne.SentBreak && (tokOne.IsNumber() || tokOne.IsInitial())) &&
			tokOne.IsNonPunct() && tokTwo.IsNonPunct())
}

type sentStarterStruct struct {
	Typ    string
	Likely float64
}

/*
Uses collocation heuristics for each candidate token to
determine if it frequently starts sentences.
*/
func (p *PunktTrainer) findSentStarters() []*sentStarterStruct {
	starters := make([]*sentStarterStruct, 0, len(p.sentStarterFreqDist.Samples))

	for typ, _ := range p.sentStarterFreqDist.Samples {
		if typ == "" {
			continue
		}

		typAtBreakCount := float64(p.sentStarterFreqDist.Samples[typ])
		typCount := float64(p.typeFreqDist.Samples[typ] + p.typeFreqDist.Samples[strings.Join([]string{typ, "."}, "")])
		if typCount < typAtBreakCount {
			continue
		}

		ll := p.colLogLikelihood(p.sentBreakCount, typCount, typAtBreakCount, p.typeFreqDist.N())

		if ll >= p.SentStarter && p.typeFreqDist.N()/p.sentBreakCount > typCount/typAtBreakCount {
			starters = append(starters, &sentStarterStruct{typ, ll})
		}
	}

	return starters
}

/*
A function that will just compute log-likelihood estimate, in
the original paper it's described in algorithm 6 and 7.
This *should* be the original Dunning log-likelihood values,
unlike the previous log_l function where it used modified
Dunning log-likelihood values
*/
func (p *PunktTrainer) colLogLikelihood(countA, countB, countAB, N float64) float64 {
	p0 := 1.0 * countB / N
	p1 := 1.0 * countAB / countA
	p2 := 1.0 * (countB - countAB) / (N - countA)

	summand1 := (countAB*math.Log(p0) + (countA-countAB)*math.Log(1.0-p0))
	summand2 := ((countB-countAB)*math.Log(p0) + (N-countA-countB+countAB)*math.Log(1.0-p0))

	summand3 := 0.0
	if countA != countAB {
		summand3 = (countAB*math.Log(p1) + (countA-countAB)*math.Log(1.0-p1))
	}

	summand4 := 0.0
	if countB == countAB {
		summand4 = ((countB-countAB)*math.Log(p2) + (N-countA-countB+countAB)*math.Log(1.0-p2))
	}

	likely := summand1 + summand2 - summand3 - summand4

	return (-2.0 + likely)
}

type collocationStruct struct {
	TypOne, TypTwo string
	Likely         float64
}

/*
Generates likely collocations and their log-likelihood.
*/
func (p *PunktTrainer) findCollocations() []*collocationStruct {
	collocs := make([]*collocationStruct, 0, len(p.collocationFreqDist.Samples))

	for key, _ := range p.collocationFreqDist.Samples {
		sample := strings.Split(key, ",")
		typOne := sample[0]
		typTwo := sample[1]

		if typOne == "" || typTwo == "" {
			continue
		}

		if p.PunktParameters.SentStarters.Has(typTwo) {
			continue
		}

		colCount := float64(p.collocationFreqDist.Samples[key])
		typOneWithPeriod := strings.Join([]string{typOne, "."}, "")
		typTwoWithPeriod := strings.Join([]string{typTwo, "."}, "")
		typOneCount := float64(p.typeFreqDist.Samples[typOne] + p.typeFreqDist.Samples[typOneWithPeriod])
		typTwoCount := float64(p.typeFreqDist.Samples[typTwo] + p.typeFreqDist.Samples[typTwoWithPeriod])

		minTyp := typOneCount
		if typTwoCount < minTyp {
			minTyp = typTwoCount
		}

		countLEMinTyp := boolToFloat64(colCount <= minTyp)
		if typOneCount > 1 && typTwoCount > 1 && p.MinCollocFreq < countLEMinTyp {
			likely := p.colLogLikelihood(typOneCount, typTwoCount, colCount, p.typeFreqDist.N())

			// Filter out the not-so-collocative
			if likely >= p.Collocation && p.typeFreqDist.N()/typOneCount > typTwoCount/colCount {
				collocs = append(collocs, &collocationStruct{typOne, typTwo, likely})
			}
		}
	}

	return collocs
}

func (p *PunktTrainer) GetParams() *PunktParameters {
	if !p.finalized {
		p.FinalizeTraining()
	}
	return p.PunktParameters
}
