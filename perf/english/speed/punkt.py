# -*- coding: utf-8 -*-
import sys
import timeit
from statistics import mean

import nltk.data

if __name__ == '__main__':
    if len(sys.argv) < 2:
        raise IOError("NLTK requires a data file, please supply the file + location as the first argument")

    if len(sys.argv) < 3:
        raise IOError("Please supply brown data file")

    fdata = None
    with open(sys.argv[2], 'r') as fp:
        fdata = fp.read()

    s = """\
        import nltk.data
        tokenizer = nltk.data.load(sys.argv[1])
        tokenier.tokenize(fdata)
    """

    print(timeit.timeit(stmt=s, number=10))
