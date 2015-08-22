import sys
import nltk

load = nltk.data.load("/home/erock/nltk_data/p2english.pickle")
sentences = load.tokenize(sys.stdin.read().strip())
for s in sentences:
    print(s)
    print("----")
