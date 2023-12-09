package barmband

type BarmbandId = [4]byte

type Barmband struct {
	Id BarmbandId
	// how often this band has found their partner (regardless which one reported)
	FoundPairs int
}

func NewBarmband(id BarmbandId) Barmband {
	return Barmband{
		Id:         id,
		FoundPairs: 0,
	}
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
