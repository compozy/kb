package convert

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/user/go-devstack/internal/models"
)

// Registry matches inputs to converters in registration order.
type Registry struct {
	converters []models.Converter
}

// UnsupportedInputError reports that no registered converter accepts the input.
type UnsupportedInputError struct {
	FilePath string
	URL      string
	Ext      string
	MIMEType string
}

// Error returns a human-readable unsupported input message.
func (e *UnsupportedInputError) Error() string {
	parts := make([]string, 0, 4)
	if e.FilePath != "" {
		parts = append(parts, fmt.Sprintf("file %q", e.FilePath))
	}
	if e.URL != "" {
		parts = append(parts, fmt.Sprintf("url %q", e.URL))
	}
	if e.Ext != "" {
		parts = append(parts, fmt.Sprintf("extension %q", e.Ext))
	}
	if e.MIMEType != "" {
		parts = append(parts, fmt.Sprintf("MIME type %q", e.MIMEType))
	}
	if len(parts) == 0 {
		return "convert: no converter accepts the provided input"
	}

	return fmt.Sprintf("convert: no converter accepts %s", strings.Join(parts, ", "))
}

// NewRegistry constructs a converter registry in priority order. When no
// converters are supplied, it registers the built-in stdlib converters.
func NewRegistry(converters ...models.Converter) *Registry {
	registry := &Registry{}
	if len(converters) == 0 {
		converters = []models.Converter{
			TextConverter{},
			HTMLConverter{},
			EPUBConverter{},
			ImageConverter{},
			PDFConverter{},
			DOCXConverter{},
			PPTXConverter{},
			XLSXConverter{},
			CSVConverter{},
			JSONConverter{},
			XMLConverter{},
		}
	}

	for _, converter := range converters {
		registry.Register(converter)
	}

	return registry
}

// Register appends a converter to the registry, preserving registration order.
func (r *Registry) Register(converter models.Converter) {
	if converter == nil {
		panic("convert: cannot register a nil converter")
	}
	r.converters = append(r.converters, converter)
}

// Match returns the first registered converter that accepts the extension or
// MIME type.
func (r *Registry) Match(ext string, mimeType string) models.Converter {
	normalizedExt := normalizeExtension(ext)
	normalizedMIMEType := normalizeMIMEType(mimeType)

	for _, converter := range r.converters {
		if converter.Accepts(normalizedExt, normalizedMIMEType) {
			return converter
		}
	}

	return nil
}

// Convert selects the first matching converter and delegates the conversion.
func (r *Registry) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	ext := inputExtension(input)
	mimeType := inputMIMEType(input.Options)

	converter := r.Match(ext, mimeType)
	if converter == nil {
		return nil, &UnsupportedInputError{
			FilePath: input.FilePath,
			URL:      input.URL,
			Ext:      ext,
			MIMEType: mimeType,
		}
	}

	return converter.Convert(ctx, input)
}

func inputExtension(input models.ConvertInput) string {
	if ext := normalizeExtension(input.FilePath); ext != "" {
		return ext
	}
	if input.URL == "" {
		return ""
	}

	parsedURL, err := url.Parse(input.URL)
	if err != nil {
		return normalizeExtension(input.URL)
	}

	return normalizeExtension(parsedURL.Path)
}

func inputMIMEType(options map[string]any) string {
	if len(options) == 0 {
		return ""
	}

	for _, key := range []string{"mimeType", "mime_type", "contentType", "content_type"} {
		value, ok := options[key]
		if !ok {
			continue
		}
		mimeType, ok := value.(string)
		if ok {
			return normalizeMIMEType(mimeType)
		}
	}

	return ""
}

func normalizeExtension(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	if ext := strings.ToLower(filepath.Ext(value)); ext != "" {
		return ext
	}

	if strings.HasPrefix(value, ".") {
		return value
	}

	if !strings.ContainsAny(value, `/\`) {
		return "." + value
	}

	return value
}

func normalizeMIMEType(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	if idx := strings.IndexByte(value, ';'); idx >= 0 {
		value = strings.TrimSpace(value[:idx])
	}

	return value
}

func readInput(input models.ConvertInput) ([]byte, error) {
	if input.Reader == nil {
		return nil, errors.New("convert: reader is required")
	}

	if _, err := input.Reader.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("convert: seek input: %w", err)
	}

	data, err := io.ReadAll(input.Reader)
	if err != nil {
		return nil, fmt.Errorf("convert: read input: %w", err)
	}

	if _, err := input.Reader.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("convert: reset input: %w", err)
	}

	return data, nil
}

func firstNonEmptyLine(body string) string {
	for _, line := range strings.Split(body, "\n") {
		title := normalizeTitleLine(line)
		if title != "" {
			return title
		}
	}

	return ""
}

func normalizeTitleLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	if strings.HasPrefix(line, "#") {
		trimmed := strings.TrimLeft(line, "#")
		if trimmed != line && strings.HasPrefix(trimmed, " ") {
			line = strings.TrimSpace(trimmed)
		}
	}

	return line
}
