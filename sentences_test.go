package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/neurosnap/sentences/punkt"
)

func loadEnglishTokenizer() *punkt.SentenceTokenizer {
	b, err := Asset("data/english.json")
	if err != nil {
		panic(err)
	}

	training, err := punkt.LoadTraining(b)
	if err != nil {
		panic(err)
	}

	return punkt.NewSentenceTokenizer(training)
}

func readFile(fname string) string {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func TestEnglish(t *testing.T) {
	tokenizer := loadEnglishTokenizer()
	tokenizer.AbbrevTypes.Add("etc")
	tokenizer.AbbrevTypes.Add("al")

	actual_text := readFile("test_files/carolyn.txt")
	expected_text := readFile("test_files/carolyn_s.txt")
	expected := strings.Split(expected_text, "\n")

	sentences := tokenizer.Tokenize(actual_text)
	for index, s := range sentences {

		if s != expected[index] {
			t.Errorf("Shit is wack")
		}
	}
}
