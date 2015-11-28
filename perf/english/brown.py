# -*- coding: utf-8 -*-
import os
import sys
import xml.etree.ElementTree as ET

if __name__ == '__main__':
    if len(sys.argv) < 2:
        raise IOError("Please provide a directory for the XML files to parse")

    if len(sys.argv) < 3:
        raise IOError("Please provide an output file for text")

    all_text = ""

    path = sys.argv[1]
    for filename in os.listdir(path):
        if not filename.endswith('.xml'):
            continue

        print("Processing file: {}".format(filename))

        fname = os.path.join(path, filename)
        tree = ET.parse(fname)

        root = tree.getroot()
        text = ET.tostring(tree.getroot(), encoding='unicode', method='text')

        # TODO: make this less hacky
        text = text.replace("\n", " ")
        text = text.replace("''", "")
        text = text.replace("``", "")
        text = text.replace(" .", ".\n")
        text = text.replace(" !", "!\n")
        text = text.replace(" ?", "?\n")

        all_text = '\n'.join([all_text, text])

    output = sys.argv[2]
    with open(output, "w") as fp:
        fp.write(all_text)

        #paragraphs = root.findall("context/p")
        #for para in paragraphs:
        #    sentences = para.findall("s")
        #    for sentence in sentences:
        #        print(sentence.text)

