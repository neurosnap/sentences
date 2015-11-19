package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VERSION string

var sentencesCmd = &cobra.Command{
	Use:   "sentences",
	Short: "Sentence tokenizer",
	Long:  "A utility that will break up a blob of text into sentences.",
	Run: func(cmd *cobra.Command, args []string) {
		//var text []byte
		fmt.Print("Hi there")
		/*if len(args) > 0 {
			text, _ = ioutil.ReadFile(args[0])
		} else {
			reader := bufio.NewReader(os.Stdin)
			text, _ = ioutil.ReadAll(reader)
		}

		tokenizer, err := english.NewSentenceTokenizer(nil)
		if err != nil {
			panic(err)
		}

		sentences := tokenizer.Tokenize(string(text))
		for _, s := range sentences {
			text := strings.Join(strings.Fields(s.Text), " ")
			fmt.Print(text)
		}*/
	},
}

func main() {
	if err := sentencesCmd.Execute(); err != nil {
		fmt.Print(err)
	}
}
