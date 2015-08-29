import sys
import nltk

load = nltk.data.load("/home/erock/nltk_data/p2english.pickle")
load._params.abbrev_types.add('etc')
load._params.abbrev_types.add('al')
sentences = load.tokenize(sys.stdin.read())
for s in sentences:
    print(s)
    print('{{sentence_break}}')
