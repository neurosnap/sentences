package english

import (
	"testing"
)

func compareSentences(t *testing.T, actualText string, expected []string, test string) bool {
	actual := tokenizer.Tokenize(actualText)

	if len(actual) != len(expected) {
		t.Log(test)
		t.Logf("Actual: %v\n", actual)
		t.Errorf("Actual: %d, Expected: %d\n", len(actual), len(expected))
		t.Log("===")
		return false
	}

	for index, sent := range actual {
		if sent.Text != expected[index] {
			t.Log(test)
			t.Errorf("Actual: [%s] Expected: [%s]\n", sent.Text, expected[index])
			t.Log("===")
			return false
		}
	}

	return true
}

func TestGoldenRules(t *testing.T) {
	var actualText string
	var expected []string
	var test string

	test = "21. Parenthetical inside sentence"
	actualText = "He teaches science (He previously worked for 5 years as an engineer.) at the local University."
	expected = []string{
		"He teaches science (He previously worked for 5 years as an engineer.) at the local University.",
	}
	compareSentences(t, actualText, expected, test)
	
	test = "24. Single quotations inside sentence"
	actualText = "She turned to him, 'This is great.' she said."
	expected = []string{
		"She turned to him, 'This is great.' she said.",
	}
	compareSentences(t, actualText, expected, test)

	test = "25. Double quotations inside sentence"
	actualText = "She turned to him, \"This is great.\" she said."
	expected = []string{
		"She turned to him, \"This is great.\" she said.",
	}
	compareSentences(t, actualText, expected, test)

	test = "26. Double quotations at the end of a sentence"
	actualText = "She turned to him, \"This is great.\" She held the book out to show him."
	expected = []string{
		"She turned to him, \"This is great.\"",
		" She held the book out to show him.",
	}
	compareSentences(t, actualText, expected, test)
}