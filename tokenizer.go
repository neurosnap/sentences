package main

// A processing interface for tokenizing a string.
type Tokenizer interface {
	tokenize(s string) []string
	span_tokenize(s string) []int
	tokenize_sents(ss []string) []string
	span_tokenize_sents(ss []string) [][]int
}
