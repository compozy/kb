package generation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode/utf8"
)

// objectSchema is the top level of a request schema, which is all the output
// check needs: required keys, additionalProperties false, and the declared
// type of each top-level property (plus `items.type` for arrays).
type objectSchema struct {
	required             []string
	properties           map[string]propertySchema
	additionalProperties bool
}

type propertySchema struct {
	types     []string
	itemTypes []string
}

func parseRequest(req Request) (*objectSchema, error) {
	switch {
	case strings.TrimSpace(req.Kind) == "":
		return nil, fmt.Errorf("%w: kind is required", ErrInvalidRequest)
	case strings.TrimSpace(req.Prompt) == "":
		return nil, fmt.Errorf("%w: prompt is required", ErrInvalidRequest)
	case strings.TrimSpace(req.SchemaName) == "":
		return nil, fmt.Errorf("%w: schema name is required", ErrInvalidRequest)
	}
	schema, err := parseSchema(req.Schema)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}
	return schema, nil
}

func parseSchema(raw json.RawMessage) (*objectSchema, error) {
	var top struct {
		Type                 json.RawMessage            `json:"type"`
		Required             []string                   `json:"required"`
		Properties           map[string]json.RawMessage `json:"properties"`
		AdditionalProperties *bool                      `json:"additionalProperties"`
	}
	if err := json.Unmarshal(raw, &top); err != nil {
		return nil, fmt.Errorf("schema must be a JSON object: %w", err)
	}
	types, err := typeList(top.Type)
	if err != nil || !slices.Equal(types, []string{"object"}) {
		return nil, errors.New(`schema type must be "object"`)
	}
	schema := &objectSchema{
		required:             top.Required,
		properties:           make(map[string]propertySchema, len(top.Properties)),
		additionalProperties: top.AdditionalProperties == nil || *top.AdditionalProperties,
	}
	for name, rawProperty := range top.Properties {
		var property struct {
			Type  json.RawMessage `json:"type"`
			Items *struct {
				Type json.RawMessage `json:"type"`
			} `json:"items"`
		}
		if err := json.Unmarshal(rawProperty, &property); err != nil {
			return nil, fmt.Errorf("schema property %q: %w", name, err)
		}
		parsed := propertySchema{}
		if parsed.types, err = typeList(property.Type); err != nil {
			return nil, fmt.Errorf("schema property %q: %w", name, err)
		}
		if property.Items != nil {
			if parsed.itemTypes, err = typeList(property.Items.Type); err != nil {
				return nil, fmt.Errorf("schema property %q items: %w", name, err)
			}
		}
		schema.properties[name] = parsed
	}
	for _, name := range schema.required {
		if _, ok := schema.properties[name]; !ok && !schema.additionalProperties {
			return nil, fmt.Errorf("required property %q is not declared", name)
		}
	}
	return schema, nil
}

// typeList reads a JSON-schema "type" (a string or an array of strings).
func typeList(raw json.RawMessage) ([]string, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return nil, nil
	}
	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		return []string{single}, nil
	}
	var many []string
	if err := json.Unmarshal(raw, &many); err != nil {
		return nil, errors.New("type must be a string or an array of strings")
	}
	return many, nil
}

// validateOutput parses content as a JSON object and checks it against the
// top level of the schema. It returns the compacted JSON.
func validateOutput(content string, schema *objectSchema) (json.RawMessage, error) {
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()
	var object map[string]any
	if err := decoder.Decode(&object); err != nil || object == nil {
		return nil, errors.New("output is not a JSON object")
	}
	if decoder.More() {
		return nil, errors.New("output has trailing data")
	}
	for _, name := range schema.required {
		if _, ok := object[name]; !ok {
			return nil, fmt.Errorf("output misses required key %q", name)
		}
	}
	for name, value := range object {
		property, declared := schema.properties[name]
		if !declared {
			if !schema.additionalProperties {
				return nil, fmt.Errorf("output has undeclared key %q", name)
			}
			continue
		}
		if len(property.types) > 0 && !matchesAny(value, property.types) {
			return nil, fmt.Errorf("output key %q has the wrong type", name)
		}
		if items, isArray := value.([]any); isArray && len(property.itemTypes) > 0 {
			for _, item := range items {
				if !matchesAny(item, property.itemTypes) {
					return nil, fmt.Errorf("output key %q has an item of the wrong type", name)
				}
			}
		}
	}
	var compact bytes.Buffer
	if err := json.Compact(&compact, []byte(strings.TrimSpace(content))); err != nil {
		return nil, fmt.Errorf("compact output: %w", err)
	}
	return compact.Bytes(), nil
}

func matchesAny(value any, types []string) bool {
	return slices.ContainsFunc(types, func(kind string) bool { return matches(value, kind) })
}

func matches(value any, kind string) bool {
	switch kind {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(json.Number)
		return ok
	case "integer":
		number, ok := value.(json.Number)
		if !ok {
			return false
		}
		parsed, err := number.Float64()
		return err == nil && parsed == math.Trunc(parsed)
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		_, ok := value.([]any)
		return ok
	case "object":
		_, ok := value.(map[string]any)
		return ok
	case "null":
		return value == nil
	default:
		return true
	}
}

// ValidateLiteral checks one generated literal: non-empty, no leading or
// trailing whitespace, no line breaks, and at most maxLen runes (maxLen ≤ 0
// disables the length check).
func ValidateLiteral(s string, maxLen int) error {
	switch {
	case strings.TrimSpace(s) == "":
		return errors.New("literal is empty")
	case strings.TrimSpace(s) != s:
		return errors.New("literal has leading or trailing whitespace")
	case strings.ContainsAny(s, "\r\n"):
		return errors.New("literal spans several lines")
	case !utf8.ValidString(s):
		return errors.New("literal is not valid UTF-8")
	case maxLen > 0 && utf8.RuneCountInString(s) > maxLen:
		return fmt.Errorf("literal has %d characters, more than %d", utf8.RuneCountInString(s), maxLen)
	}
	return nil
}

// DedupeStrings trims values, drops empty ones and case-insensitive
// duplicates, and keeps the first spelling in the original order.
func DedupeStrings(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		key := strings.ToLower(value)
		if value == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return out
}
