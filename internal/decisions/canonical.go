package decisions

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CanonicalJSON encodes v as JSON with object keys sorted, numbers kept as
// written, no HTML escaping and no trailing newline, so equal values always
// hash equally. Cache keys and value hashes are computed over it.
func CanonicalJSON(v any) ([]byte, error) {
	tree, err := toTree(v)
	if err != nil {
		return nil, err
	}
	return encodeTree(tree)
}

// HashJSON returns the hex sha256 of CanonicalJSON(v).
func HashJSON(v any) (string, error) {
	data, err := CanonicalJSON(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// toTree turns any JSON-serializable value into the generic tree
// (map[string]any, []any, string, json.Number, bool, nil).
func toTree(v any) (any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var tree any
	if err := decoder.Decode(&tree); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	return tree, nil
}

func encodeTree(tree any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(tree); err != nil {
		return nil, fmt.Errorf("encode canonical json: %w", err)
	}
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}
