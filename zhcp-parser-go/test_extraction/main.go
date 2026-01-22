package main

import (
	"fmt"
	"zhcp-parser-go/internal/parsers/pdf"
)

func main() {
	extractor := pdf.NewPDFExtractor(nil)
	result, err := extractor.ExtractText("../testdata/sample_project.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Extracted text length: %d\n", len(result.Text))
	if len(result.Text) > 0 {
		fmt.Printf("First 500 chars: %.500s\n", result.Text)
	} else {
		fmt.Println("No text extracted from PDF")
	}
}
