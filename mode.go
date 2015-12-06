package at2pt

//Mode indicates output mode
type Mode int

const (
	//PLAIN outputs plain text
	PLAIN Mode = iota
	//TOKENIZED outputs tokenized words
	TOKENIZED
)
