package convert

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/user/go-devstack/internal/models"
)

const docxMIMEType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

// DOCXConverter renders DOCX documents as Markdown paragraphs and headings.
type DOCXConverter struct{}

type docxParagraph struct {
	HeadingLevel int
	Text         string
}

type docxParagraphState struct {
	builder strings.Builder
	style   string
}

// Accepts reports whether the input is DOCX content.
func (DOCXConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".docx" || normalizeMIMEType(mimeType) == docxMIMEType
}

// Convert transforms a DOCX file into Markdown while preserving paragraphs,
// heading styles, and core document metadata when present.
func (DOCXConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	archive, err := openOfficeArchive(input, "DOCX")
	if err != nil {
		return nil, err
	}

	props, err := archive.readCoreProperties("DOCX")
	if err != nil {
		return nil, err
	}

	documentXML, err := archive.readRequiredFile("DOCX", "word/document.xml")
	if err != nil {
		return nil, err
	}

	paragraphs, err := parseDOCXParagraphs(ctx, documentXML)
	if err != nil {
		return nil, fmt.Errorf("convert: invalid DOCX content: parse word/document.xml: %w", err)
	}

	markdown := renderDOCXParagraphs(paragraphs)
	metadata := officeMetadata(props)
	if markdown == "" {
		addOfficeWarning(metadata, officeNoTextWarning)
	}

	return &models.ConvertResult{
		Markdown: markdown,
		Title:    officeDocumentTitle(props, firstNonEmptyLine(markdown)),
		Metadata: metadata,
	}, nil
}

func parseDOCXParagraphs(ctx context.Context, data []byte) ([]docxParagraph, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	paragraphs := make([]docxParagraph, 0)
	var current *docxParagraphState
	inText := false

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "p":
				current = &docxParagraphState{}
			case "pStyle":
				if current != nil {
					current.style = officeAttributeValue(token.Attr, "val")
				}
			case "t":
				if current != nil {
					inText = true
				}
			case "tab":
				if current != nil {
					current.builder.WriteString("    ")
				}
			case "br", "cr":
				if current != nil {
					current.builder.WriteString("\n")
				}
			}
		case xml.CharData:
			if inText && current != nil {
				current.builder.Write([]byte(token))
			}
		case xml.EndElement:
			switch token.Name.Local {
			case "t":
				inText = false
			case "p":
				if current == nil {
					continue
				}

				text := normalizeOfficeBlockText(current.builder.String())
				if text != "" {
					paragraphs = append(paragraphs, docxParagraph{
						HeadingLevel: docxHeadingLevel(current.style),
						Text:         text,
					})
				}

				current = nil
				inText = false
			}
		}
	}

	return paragraphs, nil
}

func renderDOCXParagraphs(paragraphs []docxParagraph) string {
	if len(paragraphs) == 0 {
		return ""
	}

	blocks := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		if paragraph.Text == "" {
			continue
		}

		if paragraph.HeadingLevel > 0 {
			blocks = append(blocks, strings.Repeat("#", paragraph.HeadingLevel)+" "+paragraph.Text)
			continue
		}

		blocks = append(blocks, paragraph.Text)
	}

	return strings.Join(blocks, "\n\n")
}

func docxHeadingLevel(style string) int {
	style = strings.TrimSpace(strings.ToLower(style))
	if !strings.HasPrefix(style, "heading") {
		return 0
	}

	value, err := strconv.Atoi(strings.TrimPrefix(style, "heading"))
	if err != nil || value < 1 {
		return 0
	}
	if value > 6 {
		return 6
	}

	return value
}

func officeAttributeValue(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return strings.TrimSpace(attr.Value)
		}
	}

	return ""
}
