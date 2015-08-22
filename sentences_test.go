package main

import (
	"fmt"
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

func TestEnglishCarolyn(t *testing.T) {
	tokenizer := loadEnglishTokenizer()
	tokenizer.AbbrevTypes.Add("etc")
	tokenizer.AbbrevTypes.Add("al")

	test_files := [][]string{
		[]string{
			"test_files/carolyn.txt",
			"test_files/carolyn_s.txt",
		},
		[]string{
			"test_files/ecig.txt",
			"test_files/ecig_s.txt",
		},
		[]string{
			"test_files/foul_ball.txt",
			"test_files/foul_ball_s.txt",
		},
		[]string{
			"test_files/fbi.txt",
			"test_files/fbi_s.txt",
		},
		[]string{
			"test_files/dre.txt",
			"test_files/dre_s.txt",
		},
		[]string{
			"test_files/dr.txt",
			"test_files/dr_s.txt",
		},
		[]string{
			"test_files/quotes.txt",
			"test_files/quotes_s.txt",
		},
		[]string{
			"test_files/kiss.txt",
			"test_files/kiss_s.txt",
		},
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
