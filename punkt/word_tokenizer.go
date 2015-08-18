package punkt

import (
	"regexp"
	"strings"
)

type PairToken struct {
	First, Second string
}

var ReMultiCharPunct string = `(?:\-{2,}|\.{2,}|(?:\.\s){2,}\.)`
var endPuncts = []string{`."`, `.'`, `.”`, ","}

func WordTokenizer(text string) []*PairToken {
	words := strings.Fields(text)
	tokens := make([]*PairToken, 0, len(words))

	multi := regexp.MustCompile(ReMultiCharPunct)

	for _, word := range words {
		// Skip one letter words
		if len(word) == 1 {
			continue
		}

		first := string(word[:1])
		second := string(word[1:])

		for _, punct := range endPuncts {
			if strings.HasSuffix(word, punct) {
				chars := strings.Split(word, "")
				first = word[:len(chars)-1]
				second = word[len(chars)-1:]
				break
			}
		}

		multipunct := multi.FindStringIndex(word)
		if multipunct != nil {
			chars := strings.Split(word, "")
			if strings.HasSuffix(word, ".") && (multipunct[1] != len(word) || multipunct[0]+multipunct[1] == len(word)) {
				first = word[:len(chars)-1]
				second = "."
			} else {
				//chars := strings.Split(multipunct, "")
				if multipunct[1] == len(word) {
					first = word[:multipunct[0]]
					second = word[multipunct[0]:]
				} else {
					first = word[:multipunct[1]]
					second = word[multipunct[1]:]
				}
			}
		}
		token := &PairToken{first, second}
		tokens = append(tokens, token)
	}

	return tokens
}
