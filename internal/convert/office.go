package convert

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/models"
)

const officeNoTextWarning = "no extractable text found"

var officeSlidePattern = regexp.MustCompile(`^ppt/slides/slide([0-9]+)\.xml$`)

type officeArchive struct {
	data   []byte
	reader *zip.Reader
}

type officeCoreProperties struct {
	Title  string `xml:"title"`
	Author string `xml:"creator"`
}

func openOfficeArchive(input models.ConvertInput, format string) (*officeArchive, error) {
	data, err := readInput(input)
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("convert: invalid %s content: %w", format, err)
	}

	return &officeArchive{
		data:   data,
		reader: reader,
	}, nil
}

func (a *officeArchive) readRequiredFile(format string, name string) ([]byte, error) {
	data, err := a.readOptionalFile(format, name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("convert: invalid %s content: missing %s", format, name)
	}

	return data, nil
}

func (a *officeArchive) readOptionalFile(format string, name string) ([]byte, error) {
	for _, file := range a.reader.File {
		if file.Name != name {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("convert: invalid %s content: read %s: %w", format, name, err)
		}

		data, readErr := io.ReadAll(reader)
		closeErr := reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("convert: invalid %s content: read %s: %w", format, name, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("convert: invalid %s content: close %s: %w", format, name, closeErr)
		}

		return data, nil
	}

	return nil, nil
}

func (a *officeArchive) readCoreProperties(format string) (officeCoreProperties, error) {
	data, err := a.readOptionalFile(format, "docProps/core.xml")
	if err != nil {
		return officeCoreProperties{}, err
	}
	if data == nil {
		return officeCoreProperties{}, nil
	}

	var props officeCoreProperties
	if err := xml.Unmarshal(data, &props); err != nil {
		return officeCoreProperties{}, fmt.Errorf("convert: invalid %s content: parse docProps/core.xml: %w", format, err)
	}

	props.Title = strings.TrimSpace(props.Title)
	props.Author = strings.TrimSpace(props.Author)
	return props, nil
}

func (a *officeArchive) matchingFiles(pattern *regexp.Regexp) []string {
	matches := make([]string, 0)
	for _, file := range a.reader.File {
		if pattern.MatchString(file.Name) {
			matches = append(matches, file.Name)
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		leftIndex := officePathNumber(matches[i])
		rightIndex := officePathNumber(matches[j])
		if leftIndex == rightIndex {
			return matches[i] < matches[j]
		}
		return leftIndex < rightIndex
	})

	return matches
}

func officePathNumber(name string) int {
	base := path.Base(name)
	matches := officeSlidePattern.FindStringSubmatch(name)
	if len(matches) == 2 {
		value, err := strconv.Atoi(matches[1])
		if err == nil {
			return value
		}
	}

	value, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(base, "slide"), ".xml"))
	if err != nil {
		return 0
	}

	return value
}

func officeMetadata(props officeCoreProperties) map[string]any {
	return map[string]any{
		"title":  props.Title,
		"author": props.Author,
	}
}

func officeDocumentTitle(props officeCoreProperties, fallback string) string {
	if props.Title != "" {
		return props.Title
	}

	return strings.TrimSpace(fallback)
}

func addOfficeWarning(metadata map[string]any, warning string) {
	if warning == "" {
		return
	}

	warnings, ok := metadata["warnings"].([]string)
	if !ok {
		warnings = []string{}
	}

	metadata["warnings"] = append(warnings, warning)
}

func normalizeOfficeBlockText(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	lines := strings.Split(value, "\n")
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line == "" {
			continue
		}
		normalized = append(normalized, line)
	}

	return strings.TrimSpace(strings.Join(normalized, "\n"))
}
