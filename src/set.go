package main

import (
	"sort"
)

type Set struct {
	keys []uint32
	data map[uint32]bool
}

func (s *Set) add(id uint32) {
	_, exists := s.data[id]

	if !exists {
		s.data[id] = true;
		s.keys = append(s.keys, id)
	}
}

func (s *Set) getByIndex(index uint32) uint32 {
	return s.keys[index]
}

func (s *Set) contains(id uint32) bool {
	_, exists := s.data[id]
	return exists
}

func (s *Set) sort() {
	sort.Slice(s.keys, func(i, j int) bool { return s.keys[i] < s.keys[j] })
}

func createSet() *Set {
	s := Set{}
	s.data = make(map[uint32]bool)

	return &s
}	