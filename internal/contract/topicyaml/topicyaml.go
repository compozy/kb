// Package topicyaml edits a topic's topic.yaml at the yaml.Node level so every
// key, comment and key order that the caller does not touch survives a write.
// Both internal/topic (metadata) and internal/contract (contract, drafts,
// decisions block) write through it.
package topicyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

// FileName is the topic metadata file name.
const FileName = "topic.yaml"

// Document is a parsed topic.yaml whose root is always a mapping node.
type Document struct {
	root *yaml.Node
	body *yaml.Node
}

// Load reads path into a Document. A missing or empty file yields an empty
// mapping. A file whose root is not a mapping is an error.
func Load(path string) (*Document, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return newDocument(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	return Parse(content)
}

// Parse decodes content into a Document.
func Parse(content []byte) (*Document, error) {
	if len(bytes.TrimSpace(content)) == 0 {
		return newDocument(), nil
	}
	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("parse topic.yaml: %w", err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return newDocument(), nil
	}
	body := root.Content[0]
	if body.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("parse topic.yaml: root must be a mapping, got %s", kindName(body.Kind))
	}
	return &Document{root: &root, body: body}, nil
}

func newDocument() *Document {
	body := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	return &Document{root: &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{body}}, body: body}
}

// Root returns the top-level mapping node.
func (d *Document) Root() *yaml.Node {
	return d.body
}

// Get returns the value node for key in mapping, or nil.
func Get(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			return mapping.Content[index+1]
		}
	}
	return nil
}

// Set encodes value and stores it under key, replacing the existing value in
// place (the key node and its comments are kept) or appending the key.
func Set(mapping *yaml.Node, key string, value any) error {
	var encoded yaml.Node
	if err := encoded.Encode(value); err != nil {
		return fmt.Errorf("encode %q: %w", key, err)
	}
	SetNode(mapping, key, &encoded)
	return nil
}

// SetNode stores node under key, replacing in place or appending.
func SetNode(mapping *yaml.Node, key string, node *yaml.Node) {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			previous := mapping.Content[index+1]
			if node.LineComment == "" {
				node.LineComment = previous.LineComment
			}
			mapping.Content[index+1] = node
			return
		}
	}
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		node,
	)
}

// Delete removes key from mapping and reports whether it was present.
func Delete(mapping *yaml.Node, key string) bool {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			mapping.Content = slices.Delete(mapping.Content, index, index+2)
			return true
		}
	}
	return false
}

// Mapping returns the mapping stored under key, creating an empty one when
// the key is missing or null. A non-mapping value is an error.
func Mapping(mapping *yaml.Node, key string) (*yaml.Node, error) {
	existing := Get(mapping, key)
	if existing != nil && existing.Kind == yaml.MappingNode {
		return existing, nil
	}
	if existing != nil && (existing.Kind != yaml.ScalarNode || existing.Tag != "!!null") {
		return nil, fmt.Errorf("%s must be a mapping, got %s", key, kindName(existing.Kind))
	}
	created := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	SetNode(mapping, key, created)
	return created, nil
}

// Bytes encodes the document with two-space indentation.
func (d *Document) Bytes() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	if err := encoder.Encode(d.root); err != nil {
		return nil, fmt.Errorf("encode topic.yaml: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("encode topic.yaml: %w", err)
	}
	return buffer.Bytes(), nil
}

// Save encodes the document and writes it to path atomically.
func (d *Document) Save(path string) error {
	content, err := d.Bytes()
	if err != nil {
		return err
	}
	return WriteFileAtomic(path, content)
}

// Update loads path, applies edit and saves the result atomically.
func Update(path string, edit func(root *yaml.Node) error) error {
	document, err := Load(path)
	if err != nil {
		return err
	}
	if err := edit(document.Root()); err != nil {
		return err
	}
	return document.Save(path)
}

// WriteFileAtomic writes content to a temporary file next to path and renames
// it into place, keeping the existing file mode (0644 for new files).
func WriteFileAtomic(path string, content []byte) error {
	mode := fs.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	dir := filepath.Dir(path)
	temp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file for %q: %w", path, err)
	}
	tempPath := temp.Name()
	cleanup := func() { _ = os.Remove(tempPath) }
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		cleanup()
		return fmt.Errorf("write temp file for %q: %w", path, err)
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		cleanup()
		return fmt.Errorf("sync temp file for %q: %w", path, err)
	}
	if err := temp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close temp file for %q: %w", path, err)
	}
	if err := os.Chmod(tempPath, mode); err != nil {
		cleanup()
		return fmt.Errorf("chmod temp file for %q: %w", path, err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		cleanup()
		return fmt.Errorf("rename temp file onto %q: %w", path, err)
	}
	return nil
}

func kindName(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "document"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.MappingNode:
		return "mapping"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	default:
		return "unknown node"
	}
}
