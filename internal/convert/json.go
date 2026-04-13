package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/compozy/kb/internal/models"
)

// JSONConverter renders JSON content inside a fenced Markdown code block.
type JSONConverter struct{}

// Accepts reports whether the input is JSON.
func (JSONConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".json" ||
		normalizeMIMEType(mimeType) == "application/json" ||
		normalizeMIMEType(mimeType) == "text/json"
}

// Convert pretty-prints JSON and extracts lightweight metadata from top-level
// scalar fields.
func (JSONConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("convert: json input is empty")
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	var formatted bytes.Buffer
	if err := json.Indent(&formatted, data, "", "  "); err != nil {
		return nil, err
	}

	return &models.ConvertResult{
		Markdown: "```json\n" + formatted.String() + "\n```",
		Title:    jsonTitle(parsed),
		Metadata: jsonMetadata(parsed),
	}, nil
}

func jsonTitle(value any) string {
	object, ok := value.(map[string]any)
	if !ok {
		return ""
	}

	for _, key := range []string{"title", "name"} {
		title, ok := object[key].(string)
		if ok && strings.TrimSpace(title) != "" {
			return strings.TrimSpace(title)
		}
	}

	return ""
}

func jsonMetadata(value any) map[string]any {
	object, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	metadata := make(map[string]any)
	for key, raw := range object {
		switch typed := raw.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				metadata[key] = strings.TrimSpace(typed)
			}
		case float64, bool, nil:
			metadata[key] = typed
		}
	}

	if len(metadata) == 0 {
		return nil
	}

	return metadata
}
