package strutil

type CharacterSetType uint

const (
	Lower CharacterSetType = 1 << iota
	Upper
	Numeric
	Symbols
)

var characterSetMap = map[CharacterSetType]string{
	Lower:   "abcdefghijklmnopqrstuvwxyz",
	Upper:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	Numeric: "0123456789",
	Symbols: "-_@.",
}

type CharacterSet struct {
	AllowedTypes CharacterSetType
	Set          string
}
