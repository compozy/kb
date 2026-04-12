package convert

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/user/kb/internal/models"
)

const (
	pptxMIMEType       = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	pptxSlideSeparator = "\n\n---\n\n"
)

// PPTXConverter renders PowerPoint presentations as Markdown grouped by slide.
type PPTXConverter struct{}

type pptxSlide struct {
	Number     int
	Paragraphs []string
}

// Accepts reports whether the input is PPTX content.
func (PPTXConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".pptx" || normalizeMIMEType(mimeType) == pptxMIMEType
}

// Convert transforms a PPTX file into Markdown sections grouped by slide and
// surfaces core document metadata when present.
func (PPTXConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	archive, err := openOfficeArchive(input, "PPTX")
	if err != nil {
		return nil, err
	}

	props, err := archive.readCoreProperties("PPTX")
	if err != nil {
		return nil, err
	}

	slideFiles := archive.matchingFiles(officeSlidePattern)
	if len(slideFiles) == 0 {
		return nil, fmt.Errorf("convert: invalid PPTX content: missing ppt/slides/slide*.xml")
	}

	slides := make([]pptxSlide, 0, len(slideFiles))
	firstText := ""
	hasText := false

	for _, slideFile := range slideFiles {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		slideXML, err := archive.readRequiredFile("PPTX", slideFile)
		if err != nil {
			return nil, err
		}

		paragraphs, err := parsePPTXSlideParagraphs(ctx, slideXML)
		if err != nil {
			return nil, fmt.Errorf("convert: invalid PPTX content: parse %s: %w", slideFile, err)
		}

		if !hasText && len(paragraphs) > 0 {
			firstText = paragraphs[0]
			hasText = true
		}

		slides = append(slides, pptxSlide{
			Number:     officePathNumber(slideFile),
			Paragraphs: paragraphs,
		})
	}

	metadata := officeMetadata(props)
	metadata["slideCount"] = len(slides)
	if !hasText {
		addOfficeWarning(metadata, officeNoTextWarning)
	}

	return &models.ConvertResult{
		Markdown: renderPPTXSlides(slides),
		Title:    officeDocumentTitle(props, firstText),
		Metadata: metadata,
	}, nil
}

func parsePPTXSlideParagraphs(ctx context.Context, data []byte) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	paragraphs := make([]string, 0)
	inTextBody := 0
	inText := false
	var current *strings.Builder

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
			case "txBody":
				inTextBody++
			case "p":
				if inTextBody > 0 {
					current = &strings.Builder{}
				}
			case "t":
				if current != nil {
					inText = true
				}
			case "br":
				if current != nil {
					current.WriteString("\n")
				}
			case "tab":
				if current != nil {
					current.WriteString("    ")
				}
			}
		case xml.CharData:
			if inText && current != nil {
				current.Write([]byte(token))
			}
		case xml.EndElement:
			switch token.Name.Local {
			case "t":
				inText = false
			case "p":
				if current == nil {
					continue
				}

				text := normalizeOfficeBlockText(current.String())
				if text != "" {
					paragraphs = append(paragraphs, text)
				}

				current = nil
				inText = false
			case "txBody":
				if inTextBody > 0 {
					inTextBody--
				}
			}
		}
	}

	return paragraphs, nil
}

func renderPPTXSlides(slides []pptxSlide) string {
	if len(slides) == 0 {
		return ""
	}

	sections := make([]string, 0, len(slides))
	for _, slide := range slides {
		var builder strings.Builder
		fmt.Fprintf(&builder, "## Slide %d", slide.Number)

		if len(slide.Paragraphs) > 0 {
			builder.WriteString("\n\n")
			builder.WriteString(strings.Join(slide.Paragraphs, "\n\n"))
		}

		sections = append(sections, builder.String())
	}

	return strings.Join(sections, pptxSlideSeparator)
}
