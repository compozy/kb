package convert

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestPDFConverterAcceptsPDFOnly(t *testing.T) {
	t.Parallel()

	converter := PDFConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".pdf", want: true},
		{mimeType: "application/pdf", want: true},
		{ext: ".txt", want: false},
		{mimeType: "text/plain", want: false},
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

func TestPDFConverterConvertsSinglePagePDFToMarkdown(t *testing.T) {
	t.Parallel()

	result := convertPDFFixture(t, "sample.pdf")

	if result.Title != "Fixture Title" {
		t.Fatalf("title = %q, want Fixture Title", result.Title)
	}

	for _, want := range []string{
		"# Welcome to Kodebase",
		"This is the first paragraph. It stays readable after conversion.",
	} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
}

func TestPDFConverterConvertsMultiPagePDFWithSeparators(t *testing.T) {
	t.Parallel()

	result := convertPDFFixture(t, "multi_page.pdf")

	if !strings.Contains(result.Markdown, pdfPageSeparator) {
		t.Fatalf("markdown %q does not contain page separator %q", result.Markdown, pdfPageSeparator)
	}
	for _, want := range []string{
		"# Overview",
		"Page one carries the introduction.",
		"## Details",
		"Page two adds the follow-up notes.",
	} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
}

func TestPDFConverterExtractsMetadataAndPageCount(t *testing.T) {
	t.Parallel()

	result := convertPDFFixture(t, "sample.pdf")

	want := map[string]any{
		"title":     "Fixture Title",
		"author":    "Fixture Author",
		"pageCount": 1,
	}
	for key, expected := range want {
		if got, ok := result.Metadata[key]; !ok || !reflect.DeepEqual(got, expected) {
			t.Fatalf("metadata[%q] = %#v, want %#v", key, got, expected)
		}
	}
}

func TestPDFConverterReturnsErrorForEncryptedPDF(t *testing.T) {
	t.Parallel()

	converter := PDFConverter{}
	data := readPDFFileBytes(t, encryptPDFFixture(t, "sample.pdf"))

	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: "encrypted.pdf",
	})
	if err == nil {
		t.Fatal("expected Convert to fail for encrypted PDF")
	}
	if !strings.Contains(err.Error(), "encrypted PDF requires a password") {
		t.Fatalf("error = %q, want encrypted PDF message", err)
	}
}

func TestPDFConverterReturnsErrorForInvalidPDFContent(t *testing.T) {
	t.Parallel()

	converter := PDFConverter{}
	data := readPDFFixtureBytes(t, "invalid.pdf")

	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: "invalid.pdf",
	})
	if err == nil {
		t.Fatal("expected Convert to fail for invalid PDF")
	}
	if !strings.Contains(err.Error(), "invalid PDF content") {
		t.Fatalf("error = %q, want invalid PDF content message", err)
	}
}

func TestPDFConverterWarnsWhenNoTextIsExtractable(t *testing.T) {
	t.Parallel()

	result := convertPDFFixture(t, "blank.pdf")

	if result.Markdown != "" {
		t.Fatalf("markdown = %q, want empty", result.Markdown)
	}
	if got, ok := result.Metadata["pageCount"]; !ok || !reflect.DeepEqual(got, 1) {
		t.Fatalf("pageCount = %#v, want 1", got)
	}
	warnings, ok := result.Metadata["warnings"].([]string)
	if !ok {
		t.Fatalf("warnings = %#v, want []string", result.Metadata["warnings"])
	}
	if len(warnings) != 1 || warnings[0] != "no extractable text found" {
		t.Fatalf("warnings = %#v, want no extractable text warning", warnings)
	}
}

func TestPDFConverterHonorsContextCancellation(t *testing.T) {
	t.Parallel()

	converter := PDFConverter{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := converter.Convert(ctx, models.ConvertInput{
		Reader:   bytes.NewReader(readPDFFixtureBytes(t, "sample.pdf")),
		FilePath: "sample.pdf",
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
}

func TestPDFWrapErrorClassifiesEncryptedInputs(t *testing.T) {
	t.Parallel()

	err := wrapPDFError("read PDF", pdfcpu.ErrWrongPassword)
	if !strings.Contains(err.Error(), "encrypted PDF requires a password") {
		t.Fatalf("wrapPDFError with wrong password = %q", err)
	}

	err = wrapPDFError("inspect PDF", errors.New("pdfcpu: this file is encrypted"))
	if !strings.Contains(err.Error(), "encrypted PDF requires a password") {
		t.Fatalf("wrapPDFError with encrypted message = %q", err)
	}

	err = wrapPDFError("inspect PDF", errors.New("boom"))
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("wrapPDFError preserved message = %q", err)
	}
}

func TestPDFHelpersSpacingAndGapClassification(t *testing.T) {
	t.Parallel()

	if !needsPDFSpace("Hello", "world") {
		t.Fatal("expected needsPDFSpace to insert a space between words")
	}
	if needsPDFSpace("Hello(", "world") {
		t.Fatal("did not expect needsPDFSpace after opening punctuation")
	}
	if got := classifyPDFGap(0, 12); got != pdfNoBreak {
		t.Fatalf("classifyPDFGap zero = %v, want %v", got, pdfNoBreak)
	}
	if got := classifyPDFGap(-12, 12); got != pdfLineBreak {
		t.Fatalf("classifyPDFGap line = %v, want %v", got, pdfLineBreak)
	}
	if got := classifyPDFGap(-24, 12); got != pdfParagraphBreak {
		t.Fatalf("classifyPDFGap paragraph = %v, want %v", got, pdfParagraphBreak)
	}
}

func TestPDFTokenParsingAndArrayExtraction(t *testing.T) {
	t.Parallel()

	content := "[(Body ) 120 <6c696e65>] TJ"

	token, next, err := nextPDFToken(content, 0)
	if err != nil {
		t.Fatalf("nextPDFToken array returned error: %v", err)
	}
	if token.kind != pdfTokenArray {
		t.Fatalf("array token kind = %v, want %v", token.kind, pdfTokenArray)
	}
	if text, ok := lastPDFArrayText([]pdfToken{token}); !ok || text != "Body line" {
		t.Fatalf("lastPDFArrayText = %q, %t, want Body line, true", text, ok)
	}

	token, _, err = nextPDFToken(content, next)
	if err != nil {
		t.Fatalf("nextPDFToken operator returned error: %v", err)
	}
	if token.kind != pdfTokenOperator || token.text != "TJ" {
		t.Fatalf("operator token = %#v, want TJ operator", token)
	}

	token, _, err = nextPDFToken("<< /MCID 0 >> /Tag", 0)
	if err != nil {
		t.Fatalf("nextPDFToken dictionary returned error: %v", err)
	}
	if token.kind != pdfTokenOther {
		t.Fatalf("dictionary token kind = %v, want %v", token.kind, pdfTokenOther)
	}

	token, _, err = nextPDFToken("(Title\\) Section)", 0)
	if err != nil {
		t.Fatalf("nextPDFToken string returned error: %v", err)
	}
	if token.text != "Title) Section" {
		t.Fatalf("string token text = %q, want Title) Section", token.text)
	}
}

func TestExtractPDFLinesSupportsTextOperators(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		"% comment",
		"BT",
		"/F1 20 Tf",
		"1 0 0 1 72 720 Tm",
		"(Title\\) Section) Tj",
		"/F1 12 Tf",
		"1 0 0 1 72 684 Tm",
		"[(Body ) 120 <6c696e65>] TJ",
		"T*",
		"(Next line) Tj",
		"ET",
	}, "\n")

	lines, err := extractPDFLines(content)
	if err != nil {
		t.Fatalf("extractPDFLines returned error: %v", err)
	}

	want := []pdfTextLine{
		{Text: "Title) Section", FontSize: 20, BreakAfter: pdfParagraphBreak},
		{Text: "Body line", FontSize: 12, BreakAfter: pdfLineBreak},
		{Text: "Next line", FontSize: 12, BreakAfter: pdfParagraphBreak},
	}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("lines = %#v, want %#v", lines, want)
	}
}

func TestParsePDFStringAndHexErrors(t *testing.T) {
	t.Parallel()

	if _, _, err := parsePDFStringLiteral("(unterminated", 0); err == nil {
		t.Fatal("expected unterminated string literal error")
	}
	if _, _, err := parsePDFHexLiteral("<4142", 0); err == nil {
		t.Fatal("expected unterminated hex literal error")
	}
	if _, err := skipPDFDictionary("<< /A << /B 1 >>", 0); err == nil {
		t.Fatal("expected unterminated dictionary error")
	}
	if _, _, err := parsePDFArray("[(alpha)", 0); err == nil {
		t.Fatal("expected unterminated array error")
	}
}

func convertPDFFixture(t *testing.T, fixtureName string) *models.ConvertResult {
	t.Helper()

	converter := PDFConverter{}
	data := readPDFFixtureBytes(t, fixtureName)

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: fixtureName,
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	return result
}

func readPDFFixtureBytes(t *testing.T, fixtureName string) []byte {
	t.Helper()

	path := filepath.Join("testdata", fixtureName)
	return readPDFFileBytes(t, path)
}

func readPDFFileBytes(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}

	return data
}

func encryptPDFFixture(t *testing.T, fixtureName string) string {
	t.Helper()

	tmpDir := t.TempDir()
	inputPath := filepath.Join("testdata", fixtureName)
	outputPath := filepath.Join(tmpDir, "encrypted.pdf")

	conf := pdfmodel.NewAESConfiguration("userpw", "ownerpw", 256)
	if err := pdfapi.EncryptFile(inputPath, outputPath, conf); err != nil {
		t.Fatalf("EncryptFile(%q) returned error: %v", inputPath, err)
	}

	return outputPath
}
