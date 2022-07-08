package gredis

type Z struct {
	Score  float64
	Member interface{}
}

type ZStore struct {
	Keys      []string
	Weights   []float64
	Aggregate string
}

type ZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}
