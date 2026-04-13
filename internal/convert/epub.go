package convert

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/compozy/kb/internal/models"
)

const epubMIMEType = "application/epub+zip"

// EPUBConverter renders EPUB archives into Markdown by walking the package
// spine order and converting XHTML chapters through the shared HTML pipeline.
type EPUBConverter struct{}

// Accepts reports whether the input is EPUB content.
func (EPUBConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".epub" || normalizeMIMEType(mimeType) == epubMIMEType
}

// Convert transforms an EPUB archive into concatenated Markdown chapters while
// surfacing package metadata from the OPF document.
func (EPUBConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}

	archive, err := openEPUBArchive(data)
	if err != nil {
		return nil, err
	}

	rootfilePath, err := archive.rootfilePath()
	if err != nil {
		return nil, err
	}

	opfData, err := archive.readRequiredFile(rootfilePath)
	if err != nil {
		return nil, err
	}

	pkg, err := parseEPUBPackage(opfData)
	if err != nil {
		return nil, fmt.Errorf("convert: invalid EPUB content: parse %s: %w", rootfilePath, err)
	}

	chapters, err := archive.markdownChapters(ctx, path.Dir(rootfilePath), pkg)
	if err != nil {
		return nil, err
	}

	markdown := strings.Join(chapters, "\n\n")
	metadata := map[string]any{
		"title":    pkg.Metadata.Title,
		"author":   pkg.Metadata.Creator,
		"language": pkg.Metadata.Language,
	}
	if len(chapters) > 0 {
		metadata["chapterCount"] = len(chapters)
	}
	if markdown == "" {
		metadata["warnings"] = []string{officeNoTextWarning}
	}

	title := strings.TrimSpace(pkg.Metadata.Title)
	if title == "" {
		title = firstNonEmptyLine(markdown)
	}

	return &models.ConvertResult{
		Markdown: markdown,
		Title:    title,
		Metadata: metadata,
	}, nil
}

type epubArchive struct {
	reader *zip.Reader
}

func openEPUBArchive(data []byte) (*epubArchive, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("convert: invalid EPUB content: %w", err)
	}

	return &epubArchive{reader: reader}, nil
}

func (a *epubArchive) rootfilePath() (string, error) {
	data, err := a.readRequiredFile("META-INF/container.xml")
	if err != nil {
		return "", err
	}

	var container epubContainer
	if err := xml.Unmarshal(data, &container); err != nil {
		return "", fmt.Errorf("convert: invalid EPUB content: parse META-INF/container.xml: %w", err)
	}
	if len(container.Rootfiles) == 0 {
		return "", fmt.Errorf("convert: invalid EPUB content: missing rootfile in META-INF/container.xml")
	}

	rootfilePath := strings.TrimSpace(container.Rootfiles[0].FullPath)
	if rootfilePath == "" {
		return "", fmt.Errorf("convert: invalid EPUB content: missing rootfile path in META-INF/container.xml")
	}

	return path.Clean(rootfilePath), nil
}

func (a *epubArchive) readRequiredFile(name string) ([]byte, error) {
	data, err := a.readOptionalFile(name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("convert: invalid EPUB content: missing %s", name)
	}

	return data, nil
}

func (a *epubArchive) readOptionalFile(name string) ([]byte, error) {
	cleanName := path.Clean(name)
	for _, file := range a.reader.File {
		if path.Clean(file.Name) != cleanName {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("convert: invalid EPUB content: read %s: %w", cleanName, err)
		}

		data, readErr := io.ReadAll(reader)
		closeErr := reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("convert: invalid EPUB content: read %s: %w", cleanName, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("convert: invalid EPUB content: close %s: %w", cleanName, closeErr)
		}

		return data, nil
	}

	return nil, nil
}

func (a *epubArchive) markdownChapters(ctx context.Context, rootDir string, pkg epubPackage) ([]string, error) {
	manifest := make(map[string]epubManifestItem, len(pkg.Manifest.Items))
	for _, item := range pkg.Manifest.Items {
		manifest[strings.TrimSpace(item.ID)] = item
	}

	chapters := make([]string, 0, len(pkg.Spine.Itemrefs))
	for _, itemref := range pkg.Spine.Itemrefs {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		item, ok := manifest[strings.TrimSpace(itemref.IDRef)]
		if !ok {
			return nil, fmt.Errorf("convert: invalid EPUB content: missing manifest item %q", itemref.IDRef)
		}

		chapterPath := path.Clean(path.Join(rootDir, item.Href))
		chapterData, err := a.readRequiredFile(chapterPath)
		if err != nil {
			return nil, err
		}

		markdown, err := HTMLToMarkdown(string(chapterData))
		if err != nil {
			return nil, fmt.Errorf("convert: convert EPUB chapter %s: %w", chapterPath, err)
		}

		markdown = strings.TrimSpace(markdown)
		if markdown == "" {
			continue
		}

		chapters = append(chapters, markdown)
	}

	return chapters, nil
}

type epubContainer struct {
	Rootfiles []epubRootfile `xml:"rootfiles>rootfile"`
}

type epubRootfile struct {
	FullPath string `xml:"full-path,attr"`
}

type epubPackage struct {
	Metadata epubPackageMetadata `xml:"metadata"`
	Manifest epubManifest        `xml:"manifest"`
	Spine    epubSpine           `xml:"spine"`
}

type epubPackageMetadata struct {
	Title    string `xml:"title"`
	Creator  string `xml:"creator"`
	Language string `xml:"language"`
}

type epubManifest struct {
	Items []epubManifestItem `xml:"item"`
}

type epubManifestItem struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type epubSpine struct {
	Itemrefs []epubSpineItemref `xml:"itemref"`
}

type epubSpineItemref struct {
	IDRef string `xml:"idref,attr"`
}

func parseEPUBPackage(data []byte) (epubPackage, error) {
	var pkg epubPackage
	if err := xml.Unmarshal(data, &pkg); err != nil {
		return epubPackage{}, err
	}

	pkg.Metadata.Title = strings.TrimSpace(pkg.Metadata.Title)
	pkg.Metadata.Creator = strings.TrimSpace(pkg.Metadata.Creator)
	pkg.Metadata.Language = strings.TrimSpace(pkg.Metadata.Language)

	return pkg, nil
}
