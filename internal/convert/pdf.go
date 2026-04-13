package convert

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/models"
	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	pdftypes "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

const pdfPageSeparator = "\n\n---\n\n"

func init() {
	// pdfcpu lazily initializes a shared config directory through global state.
	// PDF conversion only needs the built-in defaults, so disabling the config
	// dir avoids racy first-use initialization under concurrent test execution.
	pdfapi.DisableConfigDir()
}

// PDFConverter renders PDF documents as Markdown using pdfcpu for metadata and
// page content extraction.
type PDFConverter struct{}

// Accepts reports whether the input is PDF content.
func (PDFConverter) Accepts(ext string, mimeType string) bool {
	if normalizeExtension(ext) == ".pdf" {
		return true
	}

	return normalizeMIMEType(mimeType) == "application/pdf"
}

// Convert transforms a PDF into Markdown while preserving page boundaries and
// surfacing document metadata when available.
func (PDFConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}
	if !looksLikePDF(data) {
		return nil, errors.New("convert: invalid PDF content")
	}

	info, err := inspectPDF(bytes.NewReader(data), input.FilePath)
	if err != nil {
		return nil, err
	}

	pages, err := extractPDFPages(ctx, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	metadata := map[string]any{
		"title":     strings.TrimSpace(info.Title),
		"author":    strings.TrimSpace(info.Author),
		"pageCount": info.PageCount,
	}

	var warnings []string
	if allPagesEmpty(pages) {
		warnings = append(warnings, "no extractable text found")
	}
	if len(warnings) > 0 {
		metadata["warnings"] = warnings
	}

	markdown := renderPDFPages(pages)
	title := strings.TrimSpace(info.Title)
	if title == "" {
		title = firstNonEmptyLine(markdown)
	}

	return &models.ConvertResult{
		Markdown: markdown,
		Title:    title,
		Metadata: metadata,
	}, nil
}

func looksLikePDF(data []byte) bool {
	limit := len(data)
	if limit > 1024 {
		limit = 1024
	}

	return bytes.Contains(data[:limit], []byte("%PDF-"))
}

func inspectPDF(reader io.ReadSeeker, filePath string) (*pdfcpu.PDFInfo, error) {
	conf := pdfmodel.NewDefaultConfiguration()
	conf.ValidationMode = pdfmodel.ValidationRelaxed

	info, err := pdfapi.PDFInfo(reader, filepath.Base(filePath), nil, false, conf)
	if err != nil {
		return nil, wrapPDFError("inspect PDF", err)
	}

	return info, nil
}

func extractPDFPages(ctx context.Context, reader io.ReadSeeker) ([]string, error) {
	conf := pdfmodel.NewDefaultConfiguration()
	conf.ValidationMode = pdfmodel.ValidationRelaxed
	conf.Cmd = pdfmodel.EXTRACTCONTENT

	pdfCtx, err := pdfapi.ReadValidateAndOptimize(reader, conf)
	if err != nil {
		return nil, wrapPDFError("read PDF", err)
	}

	pages := make([]string, 0, pdfCtx.PageCount)
	for pageNum := 1; pageNum <= pdfCtx.PageCount; pageNum++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		pageReader, err := pdfcpu.ExtractPageContent(pdfCtx, pageNum)
		if err != nil {
			return nil, wrapPDFError(fmt.Sprintf("extract page %d", pageNum), err)
		}
		if pageReader == nil {
			pages = append(pages, "")
			continue
		}

		content, err := io.ReadAll(pageReader)
		if err != nil {
			return nil, fmt.Errorf("convert: read page %d content: %w", pageNum, err)
		}

		pageMarkdown, err := extractPDFPageMarkdown(content)
		if err != nil {
			return nil, fmt.Errorf("convert: parse page %d content: %w", pageNum, err)
		}

		pages = append(pages, pageMarkdown)
	}

	return pages, nil
}

func wrapPDFError(action string, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pdfcpu.ErrWrongPassword) || isEncryptedPDFError(err) {
		return fmt.Errorf("convert: %s: encrypted PDF requires a password", action)
	}

	return fmt.Errorf("convert: %s: %w", action, err)
}

func isEncryptedPDFError(err error) bool {
	msg := strings.ToLower(err.Error())
	for _, needle := range []string{
		"this file is encrypted",
		"correct password",
		"owner password",
		"user password",
	} {
		if strings.Contains(msg, needle) {
			return true
		}
	}

	return false
}

func allPagesEmpty(pages []string) bool {
	for _, page := range pages {
		if strings.TrimSpace(page) != "" {
			return false
		}
	}

	return true
}

func renderPDFPages(pages []string) string {
	if len(pages) == 0 || allPagesEmpty(pages) {
		return ""
	}

	return strings.Join(pages, pdfPageSeparator)
}

type pdfBreakKind int

const (
	pdfNoBreak pdfBreakKind = iota
	pdfLineBreak
	pdfParagraphBreak
)

type pdfTextLine struct {
	Text       string
	FontSize   float64
	BreakAfter pdfBreakKind
}

func extractPDFPageMarkdown(content []byte) (string, error) {
	lines, err := extractPDFLines(string(content))
	if err != nil {
		return "", err
	}

	return renderPDFTextLines(lines), nil
}

func renderPDFTextLines(lines []pdfTextLine) string {
	filtered := make([]pdfTextLine, 0, len(lines))
	for _, line := range lines {
		text := normalizePDFLine(line.Text)
		if text == "" {
			continue
		}
		line.Text = text
		filtered = append(filtered, line)
	}
	if len(filtered) == 0 {
		return ""
	}

	baseline := detectPDFFontBaseline(filtered)
	blocks := make([]string, 0, len(filtered))
	var paragraph strings.Builder

	flushParagraph := func() {
		if paragraph.Len() == 0 {
			return
		}
		blocks = append(blocks, strings.TrimSpace(paragraph.String()))
		paragraph.Reset()
	}

	for _, line := range filtered {
		if level, ok := detectHeadingLevel(line, baseline); ok {
			flushParagraph()
			blocks = append(blocks, strings.Repeat("#", level)+" "+line.Text)
			continue
		}

		if paragraph.Len() > 0 {
			paragraph.WriteByte(' ')
		}
		paragraph.WriteString(line.Text)
		if line.BreakAfter == pdfParagraphBreak {
			flushParagraph()
		}
	}

	flushParagraph()

	return strings.Join(blocks, "\n\n")
}

func normalizePDFLine(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
}

func detectPDFFontBaseline(lines []pdfTextLine) float64 {
	sizes := make([]float64, 0, len(lines))
	for _, line := range lines {
		if line.FontSize > 0 {
			sizes = append(sizes, line.FontSize)
		}
	}
	if len(sizes) == 0 {
		return 12
	}

	sort.Float64s(sizes)
	return sizes[0]
}

func detectHeadingLevel(line pdfTextLine, baseline float64) (int, bool) {
	if baseline <= 0 {
		baseline = 12
	}
	if line.FontSize <= baseline || len(line.Text) > 120 {
		return 0, false
	}

	switch {
	case line.FontSize >= baseline*1.8:
		return 1, true
	case line.FontSize >= baseline*1.5:
		return 2, true
	case line.FontSize >= baseline*1.25:
		return 3, true
	default:
		return 0, false
	}
}

type pdfTextState struct {
	lines       []pdfTextLine
	currentLine strings.Builder
	activeFont  float64
	lineFont    float64
	lastTMY     float64
	hasTMY      bool
}

func (s *pdfTextState) appendText(text string) {
	text = strings.ReplaceAll(text, "\x00", "")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	if strings.TrimSpace(text) == "" {
		return
	}

	if s.currentLine.Len() > 0 && needsPDFSpace(s.currentLine.String(), text) {
		s.currentLine.WriteByte(' ')
	}
	if s.activeFont > s.lineFont {
		s.lineFont = s.activeFont
	}
	s.currentLine.WriteString(text)
}

func (s *pdfTextState) breakLine(kind pdfBreakKind) {
	text := normalizePDFLine(s.currentLine.String())
	if text == "" {
		if len(s.lines) > 0 && s.lines[len(s.lines)-1].BreakAfter < kind {
			s.lines[len(s.lines)-1].BreakAfter = kind
		}
		s.currentLine.Reset()
		s.lineFont = 0
		return
	}

	s.lines = append(s.lines, pdfTextLine{
		Text:       text,
		FontSize:   s.lineFont,
		BreakAfter: kind,
	})
	s.currentLine.Reset()
	s.lineFont = 0
}

func needsPDFSpace(existing string, next string) bool {
	existing = strings.TrimRight(existing, " ")
	next = strings.TrimLeft(next, " ")
	if existing == "" || next == "" {
		return false
	}

	last := rune(existing[len(existing)-1])
	first := rune(next[0])

	if strings.ContainsRune("([{/", last) || strings.ContainsRune(".,;:!?)]}/", first) {
		return false
	}

	return !strings.HasSuffix(existing, " ")
}

func extractPDFLines(content string) ([]pdfTextLine, error) {
	var (
		state    pdfTextState
		operands []pdfToken
	)

	for i := 0; i < len(content); {
		token, next, err := nextPDFToken(content, i)
		if err != nil {
			return nil, err
		}
		i = next
		if token.kind == pdfTokenEOF {
			break
		}

		if token.kind != pdfTokenOperator {
			operands = append(operands, token)
			continue
		}

		switch token.text {
		case "Tf":
			if size, ok := lastPDFFloat(operands); ok {
				state.activeFont = math.Abs(size)
			}
		case "Tj":
			if text, ok := lastPDFText(operands); ok {
				state.appendText(text)
			}
		case "TJ":
			if text, ok := lastPDFArrayText(operands); ok {
				state.appendText(text)
			}
		case "'":
			state.breakLine(pdfLineBreak)
			if text, ok := lastPDFText(operands); ok {
				state.appendText(text)
			}
		case `"`:
			state.breakLine(pdfLineBreak)
			if text, ok := lastPDFText(operands); ok {
				state.appendText(text)
			}
		case "Td", "TD":
			if gap, ok := lastPDFFloat(operands); ok {
				state.breakLine(classifyPDFGap(gap, state.lineFont))
				if state.hasTMY {
					state.lastTMY += gap
				}
			}
		case "Tm":
			if y, ok := lastPDFFloat(operands); ok {
				if state.hasTMY {
					state.breakLine(classifyPDFGap(y-state.lastTMY, state.lineFont))
				}
				state.lastTMY = y
				state.hasTMY = true
			}
		case "T*":
			state.breakLine(pdfLineBreak)
		case "ET":
			state.breakLine(pdfParagraphBreak)
		}

		operands = operands[:0]
	}

	state.breakLine(pdfNoBreak)

	return state.lines, nil
}

func classifyPDFGap(deltaY float64, fontSize float64) pdfBreakKind {
	gap := math.Abs(deltaY)
	if gap < 0.01 {
		return pdfNoBreak
	}

	baseline := fontSize
	if baseline <= 0 {
		baseline = 12
	}
	if gap >= baseline*1.5 {
		return pdfParagraphBreak
	}

	return pdfLineBreak
}

type pdfTokenKind int

const (
	pdfTokenEOF pdfTokenKind = iota
	pdfTokenOperator
	pdfTokenNumber
	pdfTokenName
	pdfTokenString
	pdfTokenHexString
	pdfTokenArray
	pdfTokenOther
)

type pdfToken struct {
	kind  pdfTokenKind
	text  string
	num   float64
	array []pdfToken
}

func nextPDFToken(content string, start int) (pdfToken, int, error) {
	i := skipPDFWhitespaceAndComments(content, start)
	if i >= len(content) {
		return pdfToken{kind: pdfTokenEOF}, i, nil
	}

	switch content[i] {
	case '(':
		raw, next, err := parsePDFStringLiteral(content, i)
		if err != nil {
			return pdfToken{}, 0, err
		}
		text, err := pdftypes.StringLiteralToString(pdftypes.StringLiteral(raw))
		if err != nil {
			return pdfToken{}, 0, err
		}
		return pdfToken{kind: pdfTokenString, text: text}, next, nil
	case '<':
		if i+1 < len(content) && content[i+1] == '<' {
			next, err := skipPDFDictionary(content, i)
			if err != nil {
				return pdfToken{}, 0, err
			}
			return pdfToken{kind: pdfTokenOther}, next, nil
		}

		raw, next, err := parsePDFHexLiteral(content, i)
		if err != nil {
			return pdfToken{}, 0, err
		}
		text, err := pdftypes.HexLiteralToString(pdftypes.HexLiteral(raw))
		if err != nil {
			return pdfToken{}, 0, err
		}
		return pdfToken{kind: pdfTokenHexString, text: text}, next, nil
	case '[':
		array, next, err := parsePDFArray(content, i)
		if err != nil {
			return pdfToken{}, 0, err
		}
		return pdfToken{kind: pdfTokenArray, array: array}, next, nil
	case '/':
		next := scanPDFWord(content, i+1)
		return pdfToken{kind: pdfTokenName, text: content[i:next]}, next, nil
	}

	next := scanPDFWord(content, i)
	word := content[i:next]
	if num, err := strconv.ParseFloat(word, 64); err == nil {
		return pdfToken{kind: pdfTokenNumber, num: num, text: word}, next, nil
	}
	if isPDFOperator(word) {
		return pdfToken{kind: pdfTokenOperator, text: word}, next, nil
	}

	return pdfToken{kind: pdfTokenOther, text: word}, next, nil
}

func skipPDFWhitespaceAndComments(content string, start int) int {
	i := start
	for i < len(content) {
		switch content[i] {
		case ' ', '\t', '\n', '\r', '\f', 0:
			i++
		case '%':
			for i < len(content) && content[i] != '\n' && content[i] != '\r' {
				i++
			}
		default:
			return i
		}
	}

	return i
}

func scanPDFWord(content string, start int) int {
	i := start
	for i < len(content) {
		if strings.ContainsRune(" \t\n\r\f()<>[]{}/%", rune(content[i])) {
			break
		}
		i++
	}
	return i
}

func parsePDFStringLiteral(content string, start int) (string, int, error) {
	var (
		builder strings.Builder
		depth   = 1
		escape  bool
	)

	for i := start + 1; i < len(content); i++ {
		ch := content[i]
		switch {
		case escape:
			builder.WriteByte(ch)
			escape = false
		case ch == '\\':
			builder.WriteByte(ch)
			escape = true
		case ch == '(':
			depth++
			builder.WriteByte(ch)
		case ch == ')':
			depth--
			if depth == 0 {
				return builder.String(), i + 1, nil
			}
			builder.WriteByte(ch)
		default:
			builder.WriteByte(ch)
		}
	}

	return "", 0, errors.New("convert: unterminated PDF string literal")
}

func parsePDFHexLiteral(content string, start int) (string, int, error) {
	end := strings.IndexByte(content[start+1:], '>')
	if end < 0 {
		return "", 0, errors.New("convert: unterminated PDF hex literal")
	}

	return content[start+1 : start+1+end], start + 2 + end, nil
}

func skipPDFDictionary(content string, start int) (int, error) {
	depth := 1
	for i := start + 2; i < len(content)-1; i++ {
		switch {
		case content[i] == '<' && content[i+1] == '<':
			depth++
			i++
		case content[i] == '>' && content[i+1] == '>':
			depth--
			i++
			if depth == 0 {
				return i + 1, nil
			}
		}
	}

	return 0, errors.New("convert: unterminated PDF dictionary")
}

func parsePDFArray(content string, start int) ([]pdfToken, int, error) {
	items := make([]pdfToken, 0, 4)
	for i := start + 1; i < len(content); {
		i = skipPDFWhitespaceAndComments(content, i)
		if i >= len(content) {
			return nil, 0, errors.New("convert: unterminated PDF array")
		}
		if content[i] == ']' {
			return items, i + 1, nil
		}

		token, next, err := nextPDFToken(content, i)
		if err != nil {
			return nil, 0, err
		}
		if token.kind == pdfTokenEOF {
			return nil, 0, errors.New("convert: unterminated PDF array")
		}
		items = append(items, token)
		i = next
	}

	return nil, 0, errors.New("convert: unterminated PDF array")
}

func isPDFOperator(word string) bool {
	switch word {
	case "BT", "ET", "Tf", "Tj", "TJ", "Td", "TD", "Tm", "T*", "'", `"`:
		return true
	default:
		return false
	}
}

func lastPDFFloat(tokens []pdfToken) (float64, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].kind == pdfTokenNumber {
			return tokens[i].num, true
		}
	}

	return 0, false
}

func lastPDFText(tokens []pdfToken) (string, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i].kind {
		case pdfTokenString, pdfTokenHexString:
			return tokens[i].text, true
		}
	}

	return "", false
}

func lastPDFArrayText(tokens []pdfToken) (string, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].kind != pdfTokenArray {
			continue
		}

		var builder strings.Builder
		for _, item := range tokens[i].array {
			switch item.kind {
			case pdfTokenString, pdfTokenHexString:
				builder.WriteString(item.text)
			}
		}

		return builder.String(), true
	}

	return "", false
}
