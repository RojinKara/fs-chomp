package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"path/filepath"
)

func IsExcluded(dir string, set *hashset.Set) bool {
	return set.Contains(dir)
}

func IsExcludedFile(dir string, set *hashset.Set) bool {
	return set.Contains(filepath.Ext(dir))
}
