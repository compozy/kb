package convert

import (
	"context"
	"strings"

	"github.com/JohannesKaufmann/dom"
	markdownconverter "github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	tableplugin "github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/compozy/kb/internal/models"
	"golang.org/x/net/html"
)

// HTMLConverter renders HTML documents as Markdown.
type HTMLConverter struct{}

// Accepts reports whether the input is HTML content.
func (HTMLConverter) Accepts(ext string, mimeType string) bool {
	switch normalizeExtension(ext) {
	case ".html", ".htm":
		return true
	}

	switch normalizeMIMEType(mimeType) {
	case "text/html", "application/xhtml+xml":
		return true
	}

	return false
}

// Convert transforms HTML input into Markdown and extracts a title from the
// document metadata or first heading.
func (HTMLConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}

	markdown, title, err := convertHTMLDocument(ctx, string(data))
	if err != nil {
		return nil, err
	}

	return &models.ConvertResult{
		Markdown: markdown,
		Title:    title,
	}, nil
}

// HTMLToMarkdown exposes the shared HTML conversion pipeline for other
// converters, including EPUB chapter conversion.
func HTMLToMarkdown(htmlContent string) (string, error) {
	markdown, _, err := convertHTMLDocument(context.Background(), htmlContent)
	return markdown, err
}

func convertHTMLDocument(ctx context.Context, htmlContent string) (string, string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return "", "", err
	}

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", "", err
	}

	stripHTMLNodes(doc, "script", "style")

	title := firstElementText(doc, "title")
	if title == "" {
		title = firstHeadingText(doc)
	}

	output, err := newHTMLMarkdownConverter().ConvertNode(doc, markdownconverter.WithContext(ctx))
	if err != nil {
		return "", "", err
	}

	if err := ctx.Err(); err != nil {
		return "", "", err
	}

	markdown := strings.ReplaceAll(string(output), "\r\n", "\n")
	markdown = strings.ReplaceAll(markdown, "\r", "\n")
	markdown = strings.TrimSpace(markdown)

	return markdown, title, nil
}

func newHTMLMarkdownConverter() *markdownconverter.Converter {
	return markdownconverter.NewConverter(
		markdownconverter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithHeadingStyle(commonmark.HeadingStyleATX),
			),
			tableplugin.NewTablePlugin(),
		),
	)
}

func stripHTMLNodes(node *html.Node, names ...string) {
	if node == nil || len(names) == 0 {
		return
	}

	targets := make(map[string]struct{}, len(names))
	for _, name := range names {
		targets[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}

	stripHTMLNodesWithTargets(node, targets)
}

func stripHTMLNodesWithTargets(node *html.Node, targets map[string]struct{}) {
	for child := node.FirstChild; child != nil; {
		next := child.NextSibling

		if child.Type == html.ElementNode {
			if _, ok := targets[strings.ToLower(child.Data)]; ok {
				node.RemoveChild(child)
				child = next
				continue
			}
		}

		stripHTMLNodesWithTargets(child, targets)
		child = next
	}
}

func firstElementText(node *html.Node, name string) string {
	if node == nil {
		return ""
	}

	if node.Type == html.ElementNode && strings.EqualFold(node.Data, name) {
		return normalizeHTMLText(textContent(node))
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if text := firstElementText(child, name); text != "" {
			return text
		}
	}

	return ""
}

func firstHeadingText(node *html.Node) string {
	if node == nil {
		return ""
	}

	if node.Type == html.ElementNode && strings.EqualFold(node.Data, "h1") {
		return normalizeHTMLText(headingText(node, true))
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if text := firstHeadingText(child); text != "" {
			return text
		}
	}

	return ""
}

func textContent(node *html.Node) string {
	if node == nil {
		return ""
	}

	if node.Type == html.TextNode {
		return node.Data
	}

	var builder strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		builder.WriteString(textContent(child))
		builder.WriteString(" ")
	}

	return builder.String()
}

func headingText(node *html.Node, isRoot bool) string {
	if node == nil {
		return ""
	}

	if node.Type == html.TextNode {
		return node.Data
	}

	if node.Type == html.ElementNode {
		name := strings.ToLower(node.Data)
		if name == "br" {
			return " "
		}
		if !isRoot && dom.NameIsBlockNode(name) {
			return ""
		}
	}

	var builder strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		builder.WriteString(headingText(child, false))
		builder.WriteString(" ")
	}

	return builder.String()
}

func normalizeHTMLText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
