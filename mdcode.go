package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <directory>\n", os.Args[0])
	}

	root := os.Args[1]
  err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if info.Name() == "mdcode.go" || filepath.Ext(path) == ".mod" || filepath.Ext(path) == ".sum" {
			return nil
		}
		if !info.IsDir() {
			printFileAsMarkdown(path)
		}
		return nil
	})

  if err != nil {
    log.Fatalf("Error walking the path %q: %v\n", root, err)
  }
}

func printFileAsMarkdown(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	fmt.Printf("## %s\n", path)
	fmt.Printf("\n%s\n\n", string(content))
	//fmt.Printf("\n```\n%s\n```\n\n", string(content))
}
