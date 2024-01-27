package main

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sourcegraph/conc/iter"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type result struct {
	File string   `json:"file"`
	Line []string `json:"line"`
}

type fileResult struct {
	Path string `json:"path"`
	File string `json:"file"`
}

func main() {
	app := fiber.New()

	app.Get("/search/:path/:folderExclude/:fileExclude/:search", func(c *fiber.Ctx) error {
		excludedFolders := strings.Split(c.Params("folderExclude"), ",")
		excludedFiles := strings.Split(c.Params("fileExclude"), ",")

		excludedFoldersHashset := hashset.New()
		for _, excludedFolder := range excludedFolders {
			excludedFoldersHashset.Add(excludedFolder)
		}

		excludedFilesHashset := hashset.New()
		for _, excludedFile := range excludedFiles {
			excludedFilesHashset.Add(excludedFile)
		}

		path, _ := url.QueryUnescape(c.Params("path"))
		search, _ := url.QueryUnescape(c.Params("search"))

		results := fileWalk(path, excludedFoldersHashset, excludedFilesHashset, search)
		return c.JSON(results)
	})

	app.Get("/tree/:path/:folderExclude/:fileExclude", func(c *fiber.Ctx) error {
		excludedFolders := strings.Split(c.Params("folderExclude"), ",")
		excludedFiles := strings.Split(c.Params("fileExclude"), ",")

		excludedFoldersHashset := hashset.New()
		for _, excludedFolder := range excludedFolders {
			excludedFoldersHashset.Add(excludedFolder)
		}

		excludedFilesHashset := hashset.New()
		for _, excludedFile := range excludedFiles {
			excludedFilesHashset.Add(excludedFile)
		}

		path, _ := url.QueryUnescape(c.Params("path"))

		results := fileTree(path, excludedFoldersHashset, excludedFilesHashset)
		return c.JSON(results)
	})

	log.Fatal(app.Listen(":6969"))
}

func fileWalk(directory string, excludedFolders *hashset.Set, excludedFiles *hashset.Set, search string) []result {
	var files []string
	var results []result
	err := filepath.WalkDir(directory,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if IsExcluded(d.Name(), excludedFolders) {
				return filepath.SkipDir
			}
			if IsExcludedFile(d.Name(), excludedFiles) {
				return nil
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
			readFile, err := os.Open(fileP)

			if err != nil {
				fmt.Println(err)
			}
			fileScanner := bufio.NewScanner(readFile)
			fileScanner.Split(bufio.ScanLines)
			isText := true
			var fileLines []string

			for fileScanner.Scan() {
				fileLines = append(fileLines, fileScanner.Text())
				if !utf8.ValidString(fileScanner.Text()) {
					isText = false
					break
				}
			}

			if isText {
				matches := fuzzy.Find(search, fileLines)

				if len(matches) > 0 {
					result := &result{
						File: fileP,
						Line: matches,
					}
					results = append(results, *result)
				}

			}

			err = readFile.Close()
			if err != nil {
				return
			}
			ext := filepath.Ext(fileP)
			switch ext {
			}
		})
	return results
}

func fileTree(directory string, excludedFolders *hashset.Set, excludedFiles *hashset.Set) []fileResult {
	var files []string
	var results []fileResult
	err := filepath.WalkDir(directory,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if IsExcluded(d.Name(), excludedFolders) {
				return filepath.SkipDir
			}
			if IsExcludedFile(d.Name(), excludedFiles) {
				return nil
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
			relFileP, _ := filepath.Rel(directory, fileP)
			baseFileP := filepath.Base(fileP)
			result := &fileResult{
				Path: relFileP,
				File: baseFileP,
			}
			results = append(results, *result)
		})
	return results
}
