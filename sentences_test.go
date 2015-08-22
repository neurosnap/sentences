package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/neurosnap/sentences/punkt"
)

func loadTokenizer(data string) *punkt.SentenceTokenizer {
	b, err := Asset(data)
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

func getFileLocation(prefix, original, expected string) []string {
	orig_text := strings.Join([]string{prefix, original}, "")
	expected_text := strings.Join([]string{prefix, expected}, "")
	return []string{orig_text, expected_text}
}

func TestEnglish(t *testing.T) {
	tokenizer := loadTokenizer("data/english.json")
	tokenizer.AbbrevTypes.Add("etc")
	tokenizer.AbbrevTypes.Add("al")

	prefix := "test_files/english/"

	test_files := [][]string{
		getFileLocation(prefix, "carolyn.txt", "carolyn_s.txt"),
		getFileLocation(prefix, "ecig.txt", "ecig_s.txt"),
		getFileLocation(prefix, "foul_ball.txt", "foul_ball_s.txt"),
		getFileLocation(prefix, "fbi.txt", "fbi_s.txt"),
		getFileLocation(prefix, "dre.txt", "dre_s.txt"),
		getFileLocation(prefix, "dr.txt", "dr_s.txt"),
		getFileLocation(prefix, "quotes.txt", "quotes_s.txt"),
		getFileLocation(prefix, "kiss.txt", "kiss_s.txt"),
	}

	for _, f := range test_files {
		actual_text := readFile(f[0])
		expected_text := readFile(f[1])
		expected := strings.Split(expected_text, "\n")

		fmt.Println("Testing ", f[0])
		sentences := tokenizer.Tokenize(actual_text)
		for index, s := range sentences {
			if s != expected[index] {
				t.Errorf("%s: Actual sentence does not match expected sentence", f[0])
			}
		}
	}
}
