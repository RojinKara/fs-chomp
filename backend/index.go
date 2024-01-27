package main

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/fiber/v2"
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
	File       string `json:"file"`
	RelPath    string `json:"relPath"`
	Line       string `json:"line"`
	LineNumber int    `json:"lineNumber"`
}

type fileResult struct {
	IsFile   bool   `json:"isFile"`
	Name     string `json:"name"`
	FullPath string `json:"fullPath"`
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

	app.Get("/images", func(c *fiber.Ctx) error {
		excludedFolders := []string{"node_modules", ".git", "Applications", "tmp"}

		excludedFoldersHashset := hashset.New()
		for _, excludedFolder := range excludedFolders {
			excludedFoldersHashset.Add(excludedFolder)
		}

		path := "/Users/bhavya/Desktop"

		results := imageWalk(path, excludedFoldersHashset)
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
				for i, line := range fileLines {
					if strings.Contains(strings.ToLower(line), strings.ToLower(search)) {
						relFilePath, _ := filepath.Rel(directory, fileP)
						results = append(results, result{
							File:       fileP,
							RelPath:    relFilePath,
							Line:       line,
							LineNumber: i + 1,
						})
					}
				}
				//matches := fuzzy.Find(search, fileLines)
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
	var results []fileResult
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if IsExcluded(info.Name(), excludedFolders) {
			return filepath.SkipDir
		}
		if IsExcludedFile(info.Name(), excludedFiles) {
			return nil
		}
		relFile, _ := filepath.Rel(directory, path)
		if relFile == "." {
			return nil
		}

		if strings.Count(path, string(os.PathSeparator)) > strings.Count(directory, string(os.PathSeparator))+1 {
			return filepath.SkipDir
		}
		temp := &fileResult{
			IsFile:   !info.IsDir(),
			Name:     info.Name(),
			FullPath: path,
		}
		results = append(results, *temp)
		return nil
	})

	if err != nil {
		panic(err)
	}
	return results
}

func imageWalk(directory string, excludedFolders *hashset.Set) []string {
	var files []string
	err := filepath.WalkDir(directory,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return filepath.SkipDir
			}

			if IsExcluded(d.Name(), excludedFolders) {
				return filepath.SkipDir
			}

			if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			if !d.IsDir() && (filepath.Ext(path) == ".png" || filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg") {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	return files
}
