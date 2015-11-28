# -*- coding: utf-8 -*-
import sys
import nltk.data

if __name__ == '__main__':
    if len(sys.argv) < 2:
        raise IOError("Please supply NLTK training data")

    if len(sys.argv) < 3:
        raise IOError("Please supply text for sentence tokenizer to test")

    tokenizer = nltk.data.load(sys.argv[1])

    actual_sentences = []
    expected_sentences = []
    with open(sys.argv[2], 'r') as fp:
        data = fp.read()
        actual_sentences = tokenizer.tokenize(data)

        fp.seek(0)
        expected_sentences = fp.readlines()

    a_len = len(actual_sentences)
    e_len = len(expected_sentences)
    perc = float((e_len / a_len) * 100)

    print("Actual Sentences: {}, Expected Sentences: {}, Percent: {}" \
         .format(a_len, e_len, perc))
