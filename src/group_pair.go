package main

type GroupPair struct {
	Key string
	Cnt int
}

type byCnt []GroupPair

func (a byCnt) Len() int {
	return len(a)
}

func (a byCnt) Less(i, j int) bool {
	if a[i].Cnt == a[j].Cnt {
		return a[i].Key < a[j].Key
	} else {
		return a[i].Cnt < a[j].Cnt
	}
}

func (a byCnt) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}