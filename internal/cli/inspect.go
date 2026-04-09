package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/user/go-devstack/internal/output"
	"github.com/user/go-devstack/internal/vault"
)

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
	return vault.ReadVaultSnapshot(resolvedVault, vault.ReadVaultOptions{})
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
	)

	return command
}

func bindInspectSharedFlags(flags *pflag.FlagSet, options *inspectSharedOptions) {
	flags.StringVar(&options.Format, "format", string(output.OutputFormatTable), "Output format (table|json|tsv)")
	flags.StringVar(&options.Vault, "vault", "", "Vault root path")
	flags.StringVar(&options.Topic, "topic", "", "Topic slug inside the vault")
}

func runInspectCommand(
	cmd *cobra.Command,
	options *inspectSharedOptions,
	callback func(context inspectContext) (inspectOutput, error),
) error {
	context, err := resolveInspectContext(options)
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

	snapshot, err := readInspectVaultSnapshot(resolvedVault)
	if err != nil {
		return inspectContext{}, fmt.Errorf("inspect: read vault snapshot: %w", err)
	}

	return inspectContext{
		Format:   format,
		Snapshot: snapshot,
	}, nil
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
