package main

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sourcegraph/conc/iter"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {

}

func fileWalk(directory string, excluded *hashset.Set) {
	var files []string

	err := filepath.WalkDir(directory,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if IsExcluded(d.Name(), excluded) {
				return filepath.SkipDir
			}

			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}

	iter.ForEach(files,
		func(file *string) {
			fileP := strings.Clone(*file)
			ext := filepath.Ext(fileP)
			switch ext {
			case ".java":
				readFile, err := os.Open(fileP)

				if err != nil {
					fmt.Println(err)
				}

				fileScanner := bufio.NewScanner(readFile)
				fileScanner.Split(bufio.ScanLines)
				var fileLines []string

				for fileScanner.Scan() {
					fileLines = append(fileLines, fileScanner.Text())
				}

				matches := fuzzy.Find("getSimpleName()", fileLines)

				if len(matches) > 0 {
					fmt.Println(fileP)
					fmt.Println(matches)
					fmt.Println("--------------------------------------------------")
				}

				err = readFile.Close()
				if err != nil {
					return
				}
			}
		})
}
