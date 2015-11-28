package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/neurosnap/sentences/english"
)

func main() {
	file, err := ioutil.ReadFile("brownv.txt")
	if err != nil {
		panic(err)
	}

	text := string(file)

	var totalTime float64 = 0.0

	iterations := 10
	for i := 0; i < iterations; i++ {
		start := time.Now()

		tokenizer, err := english.NewSentenceTokenizer(nil)
		if err != nil {
			panic(err)
		}

		tokenizer.Tokenize(text)

		elapsed := time.Since(start)

		totalTime += elapsed.Seconds()

		fmt.Println("Sentences took: ", elapsed)
	}

	fmt.Println("Sentences avg took: ", totalTime/float64(iterations))

}
