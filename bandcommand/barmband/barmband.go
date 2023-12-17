package barmband

import (
	"errors"
	"fmt"
	"strconv"
)

type BarmbandId = [4]byte

var ErrInvalidLength = errors.New("barmband id has invalid length")

func IdToString(id BarmbandId) string {
	return fmt.Sprintf("%X", id)
}

// IdFromString converts  the string "12345678" to the BarmbandId []byte{0x12, 0x34, 0x56, 0x78}
func IdFromString(s string) (BarmbandId, error) {
	if len(s) != 8 {
		return BarmbandId{}, ErrInvalidLength
	}

	bytes := make([]byte, 0, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		num, _ := strconv.ParseUint(s[i:i+2], 16, 8)
		bytes = append(bytes, byte(num))
	}
	fmt.Printf("%X\n", bytes)
	return BarmbandId(bytes), nil

}

type Barmband struct {
	Id BarmbandId
	// how often this band has found their partner (regardless which one reported)
	FoundPairs int
	WantsPair  bool
}

func NewBarmband(id BarmbandId) Barmband {
	return Barmband{
		Id:         id,
		FoundPairs: 0,
		WantsPair:  false,
	}
}

type Pair struct {
	First  BarmbandId
	Second BarmbandId
	Color  string
}

func NewPair(a BarmbandId, b BarmbandId, color string) Pair {
	return Pair{
		First:  a,
		Second: b,
		Color:  color,
	}
}
