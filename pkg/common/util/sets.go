package util

import "sort"

type Empty struct{}

// sets.String is a set of strings, implemented via map[string]struct{} for minimal memory consumption.
type String map[string]Empty

// NewSet creates a String set from a list of values.
func NewSet(items ...string) String {
	ss := String{}
	ss.Insert(items...)
	return ss
}

// Insert adds items to the set.
func (s String) Insert(items ...string) String {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Delete removes all items from the set.
func (s String) Delete(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s String) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// Len returns the size of the set.
func (s String) Len() int {
	return len(s)
}

// List returns the contents as a sorted string slice.
func (s String) List() []string {
	res := make(sortableSliceOfString, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return []string(res)
}

type sortableSliceOfString []string

func (s sortableSliceOfString) Len() int           { return len(s) }
func (s sortableSliceOfString) Less(i, j int) bool { return lessString(s[i], s[j]) }
func (s sortableSliceOfString) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func lessString(lhs, rhs string) bool {
	return lhs < rhs
}
