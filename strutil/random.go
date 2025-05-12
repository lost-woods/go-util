package strutil

import (
	"crypto/rand"
	"math/big"
)

func NewCharacterSet(types CharacterSetType) CharacterSet {
	var set string
	for k, v := range characterSetMap {
		if types&k != 0 {
			set += v
		}
	}

	return CharacterSet{
		AllowedTypes: types,
		Set:          set,
	}
}

func GenerateRandomString(length int, allowedTypes CharacterSetType) string {
	charSet := NewCharacterSet(allowedTypes)

	str := make([]byte, length)
	for i := range str {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet.Set))))
		if err != nil {
			log.Fatalf("Failed to generate a random string: %v", err)
		}
		str[i] = charSet.Set[num.Int64()]
	}
	return string(str)
}
