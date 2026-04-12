package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/user/kb/internal/output"
	"github.com/user/kb/internal/vault"
)

const inspectCodebaseRelativeRoot = "raw/codebase"

type inspectSharedOptions struct {
	Format string
	Topic  string
	Vault  string
}

type inspectOutput struct {
	Columns []string
	Data    []map[string]any
}

type inspectContext struct {
	Format   output.OutputFormat
	Snapshot vault.VaultSnapshot
}

var resolveInspectVaultQuery = vault.ResolveVaultQuery

var readInspectVaultSnapshot = func(resolvedVault vault.ResolvedVault) (vault.VaultSnapshot, error) {
	snapshot, err := vault.ReadVaultSnapshot(resolvedVault, vault.ReadVaultOptions{})
	if err != nil {
		return vault.VaultSnapshot{}, err
	}

	return normalizeInspectSnapshot(snapshot), nil
}

var inspectGetwd = os.Getwd

func newInspectCommand() *cobra.Command {
	options := &inspectSharedOptions{}

	command := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a generated knowledge vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	bindInspectSharedFlags(command.PersistentFlags(), options)
	command.AddCommand(
		newInspectSmellsCommand(options),
		newInspectDeadCodeCommand(options),
		newInspectComplexityCommand(options),
		newInspectBlastRadiusCommand(options),
		newInspectCouplingCommand(options),
		newInspectSymbolCommand(options),
		newInspectFileCommand(options),
		newInspectBacklinksCommand(options),
		newInspectDepsCommand(options),
		newInspectCircularDepsCommand(options),
	)

	return command
}

func bindInspectSharedFlags(flags *pflag.FlagSet, options *inspectSharedOptions) {
	flags.StringVar(&options.Format, "format", string(output.OutputFormatTable), "Output format (table|json|tsv)")
	flags.StringVar(&options.Topic, "topic", "", "Topic slug inside the vault")
}

func runInspectCommand(
	cmd *cobra.Command,
	options *inspectSharedOptions,
	callback func(context inspectContext) (inspectOutput, error),
) error {
	resolvedOptions := *options
	resolvedOptions.Vault = commandVaultValue(cmd, options.Vault)

	context, err := resolveInspectContext(&resolvedOptions)
	if err != nil {
		return err
	}

	rendered, err := callback(context)
	if err != nil {
		return err
	}

	_, err = cmd.OutOrStdout().Write([]byte(output.FormatOutput(output.FormatOptions{
		Format:  context.Format,
		Columns: rendered.Columns,
		Data:    rendered.Data,
	})))
	if err != nil {
		return fmt.Errorf("inspect: write output: %w", err)
	}

	return nil
}

func resolveInspectContext(options *inspectSharedOptions) (inspectContext, error) {
	format, err := parseInspectOutputFormat(options.Format)
	if err != nil {
		return inspectContext{}, err
	}

	cwd, err := inspectGetwd()
	if err != nil {
		return inspectContext{}, fmt.Errorf("inspect: resolve cwd: %w", err)
	}

	resolvedVault, err := resolveInspectVaultQuery(vault.VaultQueryOptions{
		CWD:   cwd,
		Topic: strings.TrimSpace(options.Topic),
		Vault: strings.TrimSpace(options.Vault),
	})
	if err != nil {
		return inspectContext{}, err
	}

	resolvedVault.TopicPath = filepath.Join(resolvedVault.TopicPath, filepath.FromSlash(inspectCodebaseRelativeRoot))

	snapshot, err := readInspectVaultSnapshot(resolvedVault)
	if err != nil {
		return inspectContext{}, fmt.Errorf("inspect: read vault snapshot: %w", err)
	}

	return inspectContext{
		Format:   format,
		Snapshot: snapshot,
	}, nil
}

func normalizeInspectSnapshot(snapshot vault.VaultSnapshot) vault.VaultSnapshot {
	snapshot.Symbols = prefixInspectDocumentPaths(snapshot.Symbols)
	snapshot.Files = prefixInspectDocumentPaths(snapshot.Files)
	snapshot.Directories = prefixInspectDocumentPaths(snapshot.Directories)
	return snapshot
}

func prefixInspectDocumentPaths(documents []vault.VaultDocument) []vault.VaultDocument {
	if len(documents) == 0 {
		return documents
	}

	prefixed := make([]vault.VaultDocument, len(documents))
	for index, document := range documents {
		document.RelativePath = prefixInspectRelativePath(document.RelativePath)
		prefixed[index] = document
	}

	return prefixed
}

func prefixInspectRelativePath(relativePath string) string {
	trimmed := strings.Trim(strings.TrimSpace(relativePath), "/")
	if trimmed == "" {
		return ""
	}
	if trimmed == inspectCodebaseRelativeRoot || strings.HasPrefix(trimmed, inspectCodebaseRelativeRoot+"/") {
		return trimmed
	}

	return path.Join(inspectCodebaseRelativeRoot, trimmed)
}

func parseInspectOutputFormat(value string) (output.OutputFormat, error) {
	switch output.OutputFormat(strings.TrimSpace(value)) {
	case "", output.OutputFormatTable:
		return output.OutputFormatTable, nil
	case output.OutputFormatJSON:
		return output.OutputFormatJSON, nil
	case output.OutputFormatTSV:
		return output.OutputFormatTSV, nil
	default:
		return "", fmt.Errorf(`invalid --format %q: expected one of "table", "json", "tsv"`, value)
	}
}

func isFunctionLikeDocument(document vault.VaultDocument) bool {
	switch inspectFrontmatterString(document, "symbol_kind") {
	case "function", "method":
		return true
	default:
		return false
	}
}

func inspectFrontmatterString(document vault.VaultDocument, key string) string {
	value, ok := document.Frontmatter[key]
	if !ok {
		return ""
	}

	text, ok := value.(string)
	if !ok {
		return ""
	}

	return text
}

func inspectFrontmatterStringArray(document vault.VaultDocument, key string) []string {
	value, ok := document.Frontmatter[key]
	if !ok {
		return []string{}
	}

	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []any:
		values := make([]string, 0, len(typed))
		for _, entry := range typed {
			text, ok := entry.(string)
			if ok {
				values = append(values, text)
			}
		}
		return values
	default:
		return []string{}
	}
}

func inspectFrontmatterBool(document vault.VaultDocument, key string) bool {
	value, ok := document.Frontmatter[key]
	if !ok {
		return false
	}

	boolean, ok := value.(bool)
	if ok {
		return boolean
	}

	text, ok := value.(string)
	if !ok {
		return false
	}

	parsed, err := strconv.ParseBool(strings.TrimSpace(text))
	return err == nil && parsed
}

func inspectFrontmatterInt(document vault.VaultDocument, key string) int {
	value, ok := document.Frontmatter[key]
	if !ok {
		return 0
	}

	switch typed := value.(type) {
	case int:
		return typed
	case int8:
		return int(typed)
	case int16:
		return int(typed)
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case uint:
		return int(typed)
	case uint8:
		return int(typed)
	case uint16:
		return int(typed)
	case uint32:
		return int(typed)
	case uint64:
		return int(typed)
	case float32:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err == nil {
			return parsed
		}
	}

	return 0
}

func inspectFrontmatterFloat(document vault.VaultDocument, key string) float64 {
	value, ok := document.Frontmatter[key]
	if !ok {
		return 0
	}

	switch typed := value.(type) {
	case float64:
		return typed
	case float32:
		return float64(typed)
	case int:
		return float64(typed)
	case int8:
		return float64(typed)
	case int16:
		return float64(typed)
	case int32:
		return float64(typed)
	case int64:
		return float64(typed)
	case uint:
		return float64(typed)
	case uint8:
		return float64(typed)
	case uint16:
		return float64(typed)
	case uint32:
		return float64(typed)
	case uint64:
		return float64(typed)
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err == nil {
			return parsed
		}
	}

	return 0
}

type inspectDetailEntry struct {
	Field string
	Value any
}

func createInspectDetailOutput(entries ...inspectDetailEntry) inspectOutput {
	data := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		data = append(data, map[string]any{
			"field": entry.Field,
			"value": entry.Value,
		})
	}

	return inspectOutput{
		Columns: []string{"field", "value"},
		Data:    data,
	}
}

func createInspectRelationRows(relations []vault.VaultRelation) []map[string]any {
	ordered := append([]vault.VaultRelation(nil), relations...)
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].TargetPath != ordered[j].TargetPath {
			return ordered[i].TargetPath < ordered[j].TargetPath
		}
		if ordered[i].Type != ordered[j].Type {
			return ordered[i].Type < ordered[j].Type
		}
		return ordered[i].Confidence < ordered[j].Confidence
	})

	rows := make([]map[string]any, 0, len(ordered))
	for _, relation := range ordered {
		rows = append(rows, map[string]any{
			"target_path": relation.TargetPath,
			"type":        relation.Type,
			"confidence":  relation.Confidence,
		})
	}

	return rows
}

func findInspectFileBySourcePath(snapshot vault.VaultSnapshot, sourcePath string) (vault.VaultDocument, error) {
	normalizedPath := strings.TrimSpace(sourcePath)
	for _, document := range snapshot.Files {
		if inspectFrontmatterString(document, "source_path") == normalizedPath {
			return document, nil
		}
	}

	return vault.VaultDocument{}, fmt.Errorf("no file matched %q", sourcePath)
}

func findSingleInspectSymbolMatch(snapshot vault.VaultSnapshot, query string) (vault.VaultDocument, error) {
	matches := vault.FindSymbolsByName(snapshot, query)
	if len(matches) == 0 {
		return vault.VaultDocument{}, fmt.Errorf(
			"no symbols matched %q. Use `kb inspect smells` or `kb inspect complexity` to discover candidates",
			query,
		)
	}
	if len(matches) == 1 {
		return matches[0], nil
	}

	matchedNames := make([]string, 0, len(matches))
	for _, match := range matches {
		symbolName := inspectFrontmatterString(match, "symbol_name")
		if symbolName != "" {
			matchedNames = append(matchedNames, symbolName)
		}
	}

	return vault.VaultDocument{}, fmt.Errorf(
		"multiple symbols matched %q: %s. Re-run with a more specific query",
		query,
		strings.Join(matchedNames, ", "),
	)
}

func resolveInspectEntity(snapshot vault.VaultSnapshot, query string) (vault.VaultDocument, string, error) {
	if document, err := findInspectFileBySourcePath(snapshot, query); err == nil {
		return document, "file", nil
	}

	document, err := findSingleInspectSymbolMatch(snapshot, query)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no symbols matched") {
			return vault.VaultDocument{}, "", fmt.Errorf(
				"no symbol or file matched %q. Re-run with a more specific symbol name or an exact source path",
				query,
			)
		}

		return vault.VaultDocument{}, "", err
	}

	return document, "symbol", nil
}

func inspectSectionText(document vault.VaultDocument, heading string) string {
	section := strings.TrimSpace(vault.ExtractSection(document.Body, heading))
	if section == "" {
		return ""
	}

	lines := strings.Split(section, "\n")
	if len(lines) >= 3 && strings.HasPrefix(strings.TrimSpace(lines[0]), "```") && strings.TrimSpace(lines[len(lines)-1]) == "```" {
		return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
	}

	return section
}

func inspectSymbolsForFile(snapshot vault.VaultSnapshot, sourcePath string) []string {
	type symbolEntry struct {
		Name      string
		Kind      string
		StartLine int
	}

	entries := make([]symbolEntry, 0)
	for _, document := range snapshot.Symbols {
		if inspectFrontmatterString(document, "source_path") != sourcePath {
			continue
		}

		entries = append(entries, symbolEntry{
			Name:      inspectFrontmatterString(document, "symbol_name"),
			Kind:      inspectFrontmatterString(document, "symbol_kind"),
			StartLine: inspectFrontmatterInt(document, "start_line"),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].StartLine != entries[j].StartLine {
			return entries[i].StartLine < entries[j].StartLine
		}
		if entries[i].Name != entries[j].Name {
			return entries[i].Name < entries[j].Name
		}
		return entries[i].Kind < entries[j].Kind
	})

	symbols := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Kind != "" {
			symbols = append(symbols, fmt.Sprintf("%s (%s)", entry.Name, entry.Kind))
			continue
		}

		symbols = append(symbols, entry.Name)
	}

	return symbols
}
