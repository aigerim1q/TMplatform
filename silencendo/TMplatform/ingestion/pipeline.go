package ingestion

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"silencendo/sources"
)

type IngestionPipeline struct {
	chunkingOptions ChunkingOptions
}

func NewIngestionPipeline(options *ChunkingOptions) *IngestionPipeline {
	opts := ChunkingOptions{
		ChunkSize: 1000,
		Overlap:   150,
	}

	if options != nil {
		if options.ChunkSize > 0 {
			opts.ChunkSize = options.ChunkSize
		}
		if options.Overlap >= 0 {
			opts.Overlap = options.Overlap
		}
	}

	return &IngestionPipeline{
		chunkingOptions: opts,
	}
}

func (ip *IngestionPipeline) ProcessSource(source sources.KnowledgeSource) (*IngestionResult, error) {
	var text string
	var err error

	switch source.Type {
	case "file":
		text, err = ip.extractTextFromFile(source.Path)
	case "url":
		text, err = ip.extractTextFromUrl(source.Path)
	default:
		err = errors.New("Unsupported source type: " + source.Type)
	}

	if err != nil {
		return &IngestionResult{
			Chunks:   []Chunk{},
			SourceID: source.ID,
			Status:   "error",
			Error:    err.Error(),
		}, nil
	}

	chunks := ip.chunkText(text, source.ID)
	return &IngestionResult{
		Chunks:   chunks,
		SourceID: source.ID,
		Status:   "success",
	}, nil
}

func (ip *IngestionPipeline) extractTextFromFile(filePath string) (string, error) {
	resolvedPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return "", errors.New("File does not exist: " + resolvedPath)
	}

	ext := strings.ToLower(filepath.Ext(resolvedPath))

	switch ext {
	case ".txt", ".md":
		content, err := ioutil.ReadFile(resolvedPath)
		if err != nil {
			return "", err
		}
		return string(content), nil

	case ".pdf":
		// Note: For a complete implementation, you would need a PDF parsing library
		// For now, returning a placeholder since we don't have a PDF parser in our imports
		return "PDF content would be extracted here", nil

	default:
		return "", errors.New("Unsupported file type: " + ext + ". Supported types: .txt, .md, .pdf")
	}
}

func (ip *IngestionPipeline) extractTextFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "Content from " + url + " (URL content could not be retrieved)", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Content from " + url + " (URL content could not be retrieved)", nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "Content from " + url + " (URL content could not be retrieved)", nil
	}

	// Remove script and style elements
	doc.Find("script, style, nav, footer, header").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// Extract text content
	text := doc.Find("body").Text()
	return strings.Join(strings.Fields(text), " "), nil
}

func (ip *IngestionPipeline) chunkText(text string, sourceId string) []Chunk {
	var chunks []Chunk
	chunkSize := ip.chunkingOptions.ChunkSize
	overlap := ip.chunkingOptions.Overlap

	if len(text) <= chunkSize {
		chunk := Chunk{
			ChunkID:  generateUUID(),
			SourceID: sourceId,
			Text:     text,
		}
		chunk.Metadata.Position = 0
		chunk.Metadata.Length = len(text)
		chunks = append(chunks, chunk)
		return chunks
	}

	// Split text into chunks with overlap
	position := 0
	for position < len(text) {
		endPos := position + chunkSize
		if endPos > len(text) {
			endPos = len(text)
		}

		chunkText := text[position:endPos]
		chunk := Chunk{
			ChunkID:  generateUUID(),
			SourceID: sourceId,
			Text:     chunkText,
		}
		chunk.Metadata.Position = position
		chunk.Metadata.Length = len(chunkText)
		chunks = append(chunks, chunk)

		// Move position by chunkSize minus overlap
		position += chunkSize - overlap

		// If remaining text is less than or equal to overlap, we're done
		if position >= len(text) {
			break
		}
	}

	return chunks
}

// Helper function to generate a UUID
func generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant 10
	return hex.EncodeToString(bytes)
}
