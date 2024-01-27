package main

import "github.com/emirpasic/gods/sets/hashset"

func IsExcluded(dir string, set *hashset.Set) bool {
	return set.Contains(dir)
}
