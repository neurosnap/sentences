package punkt

/*
The following constants are used to describe the orthographic
contexts in which a word can occur.  BEG=beginning, MID=middle,
UNK=unknown, UC=uppercase, LC=lowercase, NC=no case.
*/
const (
	// Beginning of a sentence with upper case.
	orthoBegUc = 1 << 1
	// Middle of a sentence with upper case.
	orthoMidUc = 1 << 2
	// Unknown position in a sentence with upper case.
	orthoUnkUc = 1 << 3
	// Beginning of a sentence with lower case.
	orthoBegLc = 1 << 4
	// Middle of a sentence with lower case.
	orthoMidLc = 1 << 5
	// Unknown position in a sentence with lower case.
	orthoUnkLc = 1 << 6
	// Occurs with upper case.
	orthoUc = orthoBegUc + orthoMidUc + orthoUnkUc
	// Occurs with lower case.
	orthoLc = orthoBegLc + orthoMidLc + orthoUnkLc
)

/*
A map from context position and first-letter case to the
appropriate orthographic context flag.
*/
var orthoMap = map[[2]string]int{
	[2]string{"initial", "upper"}:  orthoBegUc,
	[2]string{"internal", "upper"}: orthoMidUc,
	[2]string{"unknown", "upper"}:  orthoUnkUc,
	[2]string{"initial", "lower"}:  orthoBegLc,
	[2]string{"internal", "lower"}: orthoMidLc,
	[2]string{"unknown", "lower"}:  orthoUnkLc,
}
