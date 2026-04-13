package convert

import (
	"archive/zip"
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestEPUBConverterAcceptsExtensionsAndMIMETypes(t *testing.T) {
	t.Parallel()

	converter := EPUBConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".epub", want: true},
		{mimeType: epubMIMEType, want: true},
		{ext: ".zip", want: false},
		{mimeType: "application/zip", want: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.ext+tc.mimeType, func(t *testing.T) {
			t.Parallel()

			if got := converter.Accepts(tc.ext, tc.mimeType); got != tc.want {
				t.Fatalf("Accepts(%q, %q) = %t, want %t", tc.ext, tc.mimeType, got, tc.want)
			}
		})
	}
}

func TestEPUBConverterConvertsFixtureInSpineOrder(t *testing.T) {
	t.Parallel()

	result := convertEPUBFixture(t, "sample.epub")

	secondIndex := strings.Index(result.Markdown, "# Second Chapter")
	firstIndex := strings.Index(result.Markdown, "# First Chapter")
	if secondIndex < 0 || firstIndex < 0 {
		t.Fatalf("markdown = %q, want both chapter headings", result.Markdown)
	}
	if secondIndex > firstIndex {
		t.Fatalf("markdown = %q, want second chapter before first chapter per spine order", result.Markdown)
	}
}

func TestEPUBConverterExtractsMetadata(t *testing.T) {
	t.Parallel()

	result := convertEPUBFixture(t, "sample.epub")

	if result.Title != "Fixture EPUB Title" {
		t.Fatalf("title = %q, want Fixture EPUB Title", result.Title)
	}

	want := map[string]any{
		"title":        "Fixture EPUB Title",
		"author":       "Fixture EPUB Author",
		"language":     "en",
		"chapterCount": 2,
	}
	for key, expected := range want {
		if got, ok := result.Metadata[key]; !ok || got != expected {
			t.Fatalf("metadata[%q] = %#v, want %#v", key, got, expected)
		}
	}
}

func TestEPUBConverterReturnsHelpfulErrorWhenPackageDocumentMissing(t *testing.T) {
	t.Parallel()

	converter := EPUBConverter{}
	data := buildEPUBBytes(t, map[string]string{
		"mimetype":               "application/epub+zip",
		"META-INF/container.xml": sampleContainerXML("OPS/content.opf"),
	})

	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: "broken.epub",
	})
	if err == nil {
		t.Fatal("expected Convert to fail")
	}
	if !strings.Contains(err.Error(), "missing OPS/content.opf") {
		t.Fatalf("error = %q, want missing package document", err)
	}
}

func TestEPUBConverterHandlesEmptyChaptersGracefully(t *testing.T) {
	t.Parallel()

	converter := EPUBConverter{}
	data := buildEPUBBytes(t, map[string]string{
		"mimetype":               "application/epub+zip",
		"META-INF/container.xml": sampleContainerXML("OPS/content.opf"),
		"OPS/content.opf": sampleOPF([]epubTestItem{
			{id: "empty", href: "empty.xhtml"},
		}, []string{"empty"}, "Empty Fixture", "Fixture Author", "en"),
		"OPS/empty.xhtml": `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml"><head><title>Empty</title></head><body></body></html>`,
	})

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: "empty.epub",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result.Markdown != "" {
		t.Fatalf("markdown = %q, want empty", result.Markdown)
	}
	if result.Title != "Empty Fixture" {
		t.Fatalf("title = %q, want Empty Fixture", result.Title)
	}

	warnings := warningList(t, result.Metadata)
	if len(warnings) != 1 || warnings[0] != officeNoTextWarning {
		t.Fatalf("warnings = %#v, want %q", warnings, officeNoTextWarning)
	}
}

func convertEPUBFixture(t *testing.T, fixtureName string) *models.ConvertResult {
	t.Helper()

	converter := EPUBConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(readConvertFixtureBytes(t, fixtureName)),
		FilePath: fixtureName,
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	return result
}

type epubTestItem struct {
	id   string
	href string
}

func buildEPUBBytes(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var archive bytes.Buffer
	writer := zip.NewWriter(&archive)

	if mimetype, ok := files["mimetype"]; ok {
		entry, err := writer.CreateHeader(&zip.FileHeader{
			Name:   "mimetype",
			Method: zip.Store,
		})
		if err != nil {
			t.Fatalf("CreateHeader(mimetype) returned error: %v", err)
		}
		if _, err := entry.Write([]byte(mimetype)); err != nil {
			t.Fatalf("Write(mimetype) returned error: %v", err)
		}
	}

	for name, content := range files {
		if name == "mimetype" {
			continue
		}

		entry, err := writer.Create(name)
		if err != nil {
			t.Fatalf("Create(%q) returned error: %v", name, err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatalf("Write(%q) returned error: %v", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}

	return archive.Bytes()
}

func sampleContainerXML(rootfilePath string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="` + rootfilePath + `" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`
}

func sampleOPF(items []epubTestItem, spine []string, title string, author string, language string) string {
	manifest := make([]string, 0, len(items))
	for _, item := range items {
		manifest = append(manifest, `    <item id="`+item.id+`" href="`+item.href+`" media-type="application/xhtml+xml"/>`)
	}

	itemrefs := make([]string, 0, len(spine))
	for _, idref := range spine {
		itemrefs = append(itemrefs, `    <itemref idref="`+idref+`"/>`)
	}

	return strings.Join([]string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<package version="3.0" xmlns="http://www.idpf.org/2007/opf" unique-identifier="bookid">`,
		`  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">`,
		`    <dc:title>` + title + `</dc:title>`,
		`    <dc:creator>` + author + `</dc:creator>`,
		`    <dc:language>` + language + `</dc:language>`,
		`  </metadata>`,
		`  <manifest>`,
		strings.Join(manifest, "\n"),
		`  </manifest>`,
		`  <spine>`,
		strings.Join(itemrefs, "\n"),
		`  </spine>`,
		`</package>`,
	}, "\n")
}
