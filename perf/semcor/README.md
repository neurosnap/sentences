
                              SemCor 3.0 
                        =====================
                            June 13, 2008 
                     Rada Mihalcea, rada@cs.unt.edu
                       University of North Texas 



SemCor 3.0 was automatically created from SemCor 1.6 by mapping 
WordNet 1.6 to WordNet 3.0 senses.
SemCor 1.6 was created and is property of Princeton University.

Some (few) word senses from WordNet 1.6 were dropped, and therefore 
they cannot be retrieved anymore in the 3.0 database. A sense of 0 
(wnsn=0) is used to symbolize a missing sense in WordNet 3.0.

The automatic mapping was performed within the Language and Information 
Technologies lab at UNT, by Rada Mihalcea (rada@cs.unt.edu).

THIS MAPPING IS PROVIDED "AS IS" AND UNT MAKES NO REPRESENTATIONS 
OR WARRANTIES, EXPRESS OR IMPLIED.  BY WAY OF EXAMPLE, BUT NOT 
LIMITATION, UNT MAKES NO REPRESENTATIONS OR WARRANTIES OF MERCHANT-
ABILITY OR FITNESS FOR ANY PARTICULAR PURPOSE.

In agreement with the license from Princeton Univerisity, you are 
granted permission to use, copy, modify and distribute this database  
for any purpose and without fee and royalty is hereby granted, provided 
that you agree to comply with the Princeton copyright notice and 
statements, including the disclaimer, and that the same appear on ALL 
copies of the database, including modifications that you make for internal  
use or for distribution.  
Both LICENSE and README files distributed with the SemCor 1.6 package
are included in the current distribution of SemCor 3.0.

When the WordNet Semantic Concordance package is unbundled you should
have the following files and subdirectories in this directory:

	README	   	this file
	LICENSE			WordNet copyright and license agreement
	INSTALL			INSTALL file for SemCor 1.6 
	brown1		103 semantically tagged Brown Corpus files
				(all content words tagged)
	brown2		83 semantically tagged Brown Corpus files
				(all content words tagged)
	brownv		166 semantically tagged Brown Corpus files
				(only verbs tagged)

Each semantic concordance directory (brown1, brown2, brownv) contains
the following subdirectories:
	tagfiles	directory of semantically tagged files

To access the files, unzip and untar the file semcor3.0.tar.gz.

For other files, such as documentation, taglists, etc., please check 
the SemCor 1.6 distribution (available from this site).

Any questions regarding this mapping should be addressed to 
Rada Mihalcea (rada@cs.unt.edu)

SemCor files have been converted to well-formed XML, and given
an ".xml" filename extension, for use with NLTK's SemCor corpus reader.

