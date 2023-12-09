package barmband

type BarmbandId = [4]byte

type Barmband struct {
	Id BarmbandId
}

type Pair struct {
	First  BarmbandId
	Second BarmbandId
}

func NewPair(a BarmbandId, b BarmbandId) Pair {
	return Pair{
		First:  a,
		Second: b,
	}
}
