package main

type Bucket struct {
	data []IterableSet
}

func (s *Bucket) count() int {
	result := 0

	for _, b := range s.data {
		result += len(b.set.keys)
	}

	return result
}

func (s *Bucket) contains(id uint32) bool {
	for _, b := range s.data {
		_, exists := b.set.data[id]

		if exists {
			return true;
		}
	}

	return false
}

func (s *Bucket) atEnd() bool {
	for _, b := range s.data {
		if b.index >= 0 {
			return false
		}
	}

	return true
}

func (s *Bucket) getMaxAndMove() uint32 {
	var maxValue uint32 = 0
	var maxIndex int = 0
	
	for index, b := range s.data {
		if b.index >= 0 {
			if maxValue == 0 || maxValue < b.set.keys[b.index] {
				maxValue = b.set.keys[b.index]
				maxIndex = index
			}
		}
	}

	s.data[maxIndex].index--

	return maxValue
}

func createBucket(set *Set) Bucket {
	var pair IterableSet
	pair.set = set
	pair.index = len(set.keys) - 1

	var bucket Bucket
	bucket.data = append(bucket.data, pair)

	return bucket
}

func (bucket *Bucket) add(set *Set) {
	var pair IterableSet
	pair.set = set
	pair.index = len(set.keys) - 1

	bucket.data = append(bucket.data, pair)
}

type byCount []Bucket

func (a byCount) Len() int {
	return len(a)
}

func (a byCount) Less(i, j int) bool {
	return a[i].count() < a[j].count()
}

func (a byCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
