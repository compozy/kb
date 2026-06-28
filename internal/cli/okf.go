package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/models"
	kokf "github.com/compozy/kb/internal/okf"
	"github.com/compozy/kb/internal/output"
	ktopic "github.com/compozy/kb/internal/topic"
)

type okfCheckOptions struct {
	Format string
	Strict bool
}

var runOKFCheck = kokf.Check
var runOKFTopicInfo = ktopic.Info
var okfGetwd = os.Getwd

func newOKFCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "okf",
		Short: "Work with Open Knowledge Format bundles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(newOKFCheckCommand())
	return command
}

func newOKFCheckCommand() *cobra.Command {
	options := &okfCheckOptions{
		Format: string(output.OutputFormatTable),
	}
	command := &cobra.Command{
		Use:   "check <topic>",
		Short: "Check an OKF topic for conformance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOKFCheckCommand(cmd, options, args[0])
		},
	}
	command.Flags().StringVar(&options.Format, "format", string(output.OutputFormatTable), "Output format (table|json|tsv)")
	command.Flags().BoolVar(&options.Strict, "strict", false, "Promote local-standard warnings to errors")
	return command
}

func runOKFCheckCommand(cmd *cobra.Command, options *okfCheckOptions, topicSlug string) error {
	format, err := parseInspectOutputFormat(options.Format)
	if err != nil {
		return err
	}
	vaultPath, err := resolveCommandVaultPath(cmd, okfGetwd, "okf check")
	if err != nil {
		return err
	}
	cfg, err := loadCLIConfig()
	if err != nil {
		return fmt.Errorf("okf check: %w", err)
	}
	topicInfo, err := runOKFTopicInfo(vaultPath, topicSlug)
	if err != nil {
		return fmt.Errorf("okf check: %w", err)
	}
	if topicInfo.Mode != models.TopicModeOKF {
		return fmt.Errorf("okf check: topic %q is not an OKF topic", topicSlug)
	}

	issues, err := runOKFCheck(commandContext(cmd), topicInfo.RootPath, kokf.CheckOptions{
		Types:  cfg.OKF.Types,
		Strict: options.Strict,
	})
	if err != nil {
		return err
	}
	_, writeErr := cmd.OutOrStdout().Write([]byte(output.FormatOutput(output.FormatOptions{
		Format:  format,
		Columns: kokf.Columns(),
		Data:    kokf.Rows(issues),
	})))
	if writeErr != nil {
		return fmt.Errorf("okf check: write output: %w", writeErr)
	}
	if kokf.HasErrors(issues) {
		return fmt.Errorf("okf check: found %d issue(s)", len(issues))
	}
	return nil
}
