package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Characters that are candidates for sentence boundaries
var sent_end_chars = []string{".", "?", "!"}

func re_sent_end_chars(expressions []string) string {
	return regexp.QuoteMeta(strings.Join(sent_end_chars, ""))
}

// Sentence internal punctuation, which indicates an abbreviation if preceded by a period-final token
var internal_punctuation = ",:;"

// Used to realign punctuation that should be included in a sentence although it follows the period (or ?, !)
var re_boundary_realignment = regexp.MustCompile(`[\"')\]}]+?(?:\s+|(?=--)|$)`)

// Excludes some characters from starting word tokens
var re_word_start = regexp.MustCompile("[^\\(\"\\`{\\[:;&\\#\\*@\\)}\\]\\-,]")

// Characters that cannot appear within words
var re_non_word_chars = regexp.MustCompile(`(?:[?!)\";}\]\*:@\'\({\[])`)

// Hyphen and ellipsis are multi-character punctuation
var re_multi_char_punct = regexp.MustCompile(`(?:\-{2,}|\.{2,}|(?:\.\s){2,}\.)`)

var tokenize_fmt = `(
%v
|
(?=%v)\S+?
(?=
\s|
$|
%v|%v|
,(?=$|\s|%v|%v)
)
|
\S
)`

var work_tokenize_fmt = regexp.MustCompile(fmt.Sprintf(tokenize_fmt, re_multi_char_punct, re_word_start, re_non_word_chars, re_multi_char_punct, re_non_word_chars, re_multi_char_punct))

// Tokenize a string to split off punctuation other than periods
func word_tokenize(s string) []string {
	return work_tokenize_fmt.FindAllString(s, -1)
}

// Format of a regular expression to find contexts including possible
// sentence boundaries. Matches token which the possible sentence boundary
// ends, and matches the following token within a lookahead expression
var period_context_fmt = `
\S*
%v
(?=(?P<after_tok>
%v
|
\s+(?P<next_tok>\S+)
))`

var period_context = regexp.MustCompile(fmt.Sprintf(period_context_fmt, re_sent_end_chars, re_non_word_chars))

var re_non_punct = regexp.MustCompile(`[^\W\d]`)
