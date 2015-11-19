package english

import (
	"testing"
)

var tokenizer, _ = NewSentenceTokenizer(nil)

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
		if sent.Text != expected[index] {
			t.Fatalf("Actual: %s\nExpected: %s", sent.Text, expected[index])
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
		if sent.Text != expected[index] {
			t.Fatalf("Actual: %s\nExpected: %s", sent.Text, expected[index])
		}
	}

	actual_text = "An abbreviation near the end of a G.D. sentence.  J.G. Wentworth was cool."
	actual = tokenizer.Tokenize(actual_text)

	expected = []string{
		"An abbreviation near the end of a G.D. sentence.",
		"  J.G. Wentworth was cool.",
	}

	if len(actual) != len(expected) {
		t.Fatalf("Actual: %d, Expected: %d", len(actual), len(expected))
	}

	for index, sent := range actual {
		if sent.Text != expected[index] {
			t.Fatalf("Actual: %s\nExpected: %s", sent.Text, expected[index])
		}
	}
}

func TestEnglishSupervisedAbbrev(t *testing.T) {
	t.Log("Tokenizer should detect list of supervised abbreviations.")

	actual_text := "I am a Sgt. in the army.  I am a No. 1 student.  The Gov. of Michigan is a dick."
	actual := tokenizer.Tokenize(actual_text)

	expected := []string{
		"I am a Sgt. in the army.",
		"  I am a No. 1 student.",
		"  The Gov. of Michigan is a dick.",
	}

	if len(actual) != len(expected) {
		t.Fatalf("Actual: %d, Expected: %d", len(actual), len(expected))
	}

	for index, sent := range actual {
		if sent.Text != expected[index] {
			t.Fatalf("Actual: %s\nExpected: %s", sent.Text, expected[index])
		}
	}
}
