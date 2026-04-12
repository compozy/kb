package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	klint "github.com/user/kb/internal/lint"
	"github.com/user/kb/internal/output"
	"github.com/user/kb/internal/vault"
)

type lintCommandOptions struct {
	Format string
	Save   bool
	Topic  string
	Vault  string
}

var runLintEngine = klint.Lint
var saveLintEngineReport = klint.SaveReport
var resolveLintVaultQuery = vault.ResolveVaultQuery
var lintGetwd = os.Getwd
var lintNow = time.Now

func newLintCommand() *cobra.Command {
	options := &lintCommandOptions{
		Format: string(output.OutputFormatTable),
	}

	command := &cobra.Command{
		Use:   "lint [<slug>]",
		Short: "Check one topic for structural KB issues",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLintCommand(cmd, options, args)
		},
	}

	flags := command.Flags()
	flags.StringVar(&options.Format, "format", string(output.OutputFormatTable), "Output format (table|json|tsv)")
	flags.BoolVar(&options.Save, "save", false, "Write a markdown report to outputs/reports/<date>-lint.md")
	flags.StringVar(&options.Topic, "topic", "", "Topic slug inside the vault")

	return command
}

func runLintCommand(cmd *cobra.Command, options *lintCommandOptions, args []string) error {
	format, err := parseInspectOutputFormat(options.Format)
	if err != nil {
		return err
	}

	topicSlug, err := resolveLintTopicSlug(args, options.Topic)
	if err != nil {
		return err
	}

	cwd, err := lintGetwd()
	if err != nil {
		return fmt.Errorf("lint: resolve cwd: %w", err)
	}

	resolvedVault, err := resolveLintVaultQuery(vault.VaultQueryOptions{
		CWD:   cwd,
		Topic: topicSlug,
		Vault: commandVaultValue(cmd, options.Vault),
	})
	if err != nil {
		return fmt.Errorf("lint: %w", err)
	}

	issues, err := runLintEngine(resolvedVault.TopicPath)
	if err != nil {
		return fmt.Errorf("lint: %w", err)
	}

	if options.Save {
		if _, err := saveLintEngineReport(resolvedVault.TopicPath, issues, lintNow()); err != nil {
			return fmt.Errorf("lint: %w", err)
		}
	}

	_, err = cmd.OutOrStdout().Write([]byte(output.FormatOutput(output.FormatOptions{
		Format:  format,
		Columns: klint.Columns(),
		Data:    klint.Rows(issues),
	})))
	if err != nil {
		return fmt.Errorf("lint: write output: %w", err)
	}

	return nil
}

func resolveLintTopicSlug(args []string, flagValue string) (string, error) {
	positionalTopic := ""
	if len(args) > 0 {
		positionalTopic = strings.TrimSpace(args[0])
		if positionalTopic == "" {
			return "", fmt.Errorf("lint: topic slug must not be empty")
		}
	}

	flagTopic := strings.TrimSpace(flagValue)
	switch {
	case positionalTopic == "":
		return flagTopic, nil
	case flagTopic == "":
		return positionalTopic, nil
	case positionalTopic != flagTopic:
		return "", fmt.Errorf(
			"lint: received conflicting topic selectors %q and %q; use either the positional slug or --topic",
			positionalTopic,
			flagTopic,
		)
	default:
		return positionalTopic, nil
	}
}
