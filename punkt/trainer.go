package punkt

import (
	gs "github.com/neurosnap/go-sentences"
)

// Learns parameters used in Punkt sentence boundary detection
type PunktTrainer struct {
	PunktBase
	typeFreqDist         *gs.FreqDist
	numPeriodToks        int
	collocationFreqDist  *gs.FreqDist
	sentStarterFreqDist  *gs.FreqDist
	sentBreakCount       int
	finalized            bool
	Abbrev               float32
	IgnoreAbbrevPenalty  bool
	AbbrevBackoff        int
	Collocation          float64
	SentStarter          int
	IncludeAllCollocs    bool
	IncludeAbbrevCollocs bool
	MinCollocFreq        int
}

func NewPunktTrainer(trainText string, verbose bool, languageVars PunktLanguageVars, token PunktToken) *PunktTrainer {
	trainer := &PunktTrainer{
		typeFreqDist:        &gs.FreqDist{},
		collocationFreqDist: &gs.FreqDist{},
		sentStarterFreqDist: &gs.FreqDist{},
		finalized:           true,
		Abbrev:              0.3,
		AbbrevBackoff:       5,
		Collocation:         7.88,
		SentStarter:         30,
		MinCollocFreq:       1,
	}

	if trainText != "" {
		trainer.Train(trainText, verbose, true)
	}

	return trainer
}

func (p *PunktTrainer) Train(text string, verbose bool, finalize bool) {
	p.TrainTokens(p.TokenizeWords(text), verbose, finalize)
	if finalize {
		p.FinalizeTraining(verbose)
	}
}

func (p *PunktTrainer) TrainTokens(tokens []*PunktToken, verbose bool, finalize bool) {

}

func (p *PunktTrainer) FinalizeTraining(verbose bool) {}
