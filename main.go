package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-md-to-html <input-dir> <output-dir>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]

	err := processDirectory(inputDir, outputDir)
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}
}

func processDirectory(inputDir, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				return err
			}
			outPath := filepath.Join(outputDir, relPath)
			return os.MkdirAll(outPath, 0755)
		}

		if filepath.Ext(path) == ".md" {
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				return err
			}
			outPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".md")+".html")
			return convertFile(path, outPath)
		}

		return nil
	})
}

func convertFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	parser := NewParser(string(input))
	html := parser.Parse()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.WriteString(outFile, html)
	return err
}
