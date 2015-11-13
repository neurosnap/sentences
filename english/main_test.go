package english

import (
	"testing"
)

var tokenizer = NewSentenceTokenizer(nil)

func TestEnglishSmartQuotes(t *testing.T) {
	t.Log("Tokenizer should break sentences that end in smart quotes ...")

	actual_text := "Here is a quote, ”a smart one.” Will this break properly?"
	actual := tokenizer.Tokenize(actual_text)

	expected := []string{
		"Here is a quote, ”a smart one.”",
		" Will this break properly?",
	}

	if len(actual) != len(expected) {
		t.Fatalf("Actual: %d, Expected: %d", len(actual), len(expected))
	}

	for index, sent := range actual {
		if sent != expected[index] {
			t.Fatalf("Actual: %s\nExpected: %s", sent, expected[index])
		}
	}
}

func TestEnglishCustomAbbrev(t *testing.T) {
	t.Log("Tokenizer should detect custom abbreviations and not always sentence break on them.")

	actual_text := "One custom abbreviation is F.B.I.  The abbreviation, F.B.I. should properly break."
	actual := tokenizer.Tokenize(actual_text)

	expected := []string{
		"One custom abbreviation is F.B.I.",
		"  The abbreviation, F.B.I. should properly break.",
	}

	if len(actual) != len(expected) {
		t.Fatalf("Actual: %d, Expected: %d", len(actual), len(expected))
	}

	for index, sent := range actual {
		if sent != expected[index] {
			t.Fatalf("Actual: %s\nExpected: %s", sent, expected[index])
		}
	}
}
