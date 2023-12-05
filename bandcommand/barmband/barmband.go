package barmband

type BarmbandId = [4]byte

type Barmband struct {
	Id BarmbandId
}

type Pair struct {
	First  BarmbandId
	Second BarmbandId
}
