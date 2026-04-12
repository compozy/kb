package convert

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"strings"

	"github.com/user/kb/internal/models"
)

// XMLConverter extracts text content from XML documents.
type XMLConverter struct{}

// Accepts reports whether the input is XML.
func (XMLConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".xml" ||
		normalizeMIMEType(mimeType) == "application/xml" ||
		normalizeMIMEType(mimeType) == "text/xml"
}

// Convert strips tags and keeps the document's text content.
func (XMLConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("convert: xml input is empty")
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	fragments := make([]string, 0, 8)

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		charData, ok := token.(xml.CharData)
		if !ok {
			continue
		}

		text := strings.Join(strings.Fields(string(charData)), " ")
		if text != "" {
			fragments = append(fragments, text)
		}
	}

	if len(fragments) == 0 {
		return nil, errors.New("convert: xml input has no text content")
	}

	return &models.ConvertResult{
		Markdown: strings.Join(fragments, " "),
		Title:    fragments[0],
	}, nil
}
