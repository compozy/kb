package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/models"
	kokf "github.com/compozy/kb/internal/okf"
	"github.com/compozy/kb/internal/output"
	"github.com/compozy/kb/internal/session"
	ktopic "github.com/compozy/kb/internal/topic"
)

type okfCheckOptions struct {
	Format string
	Strict bool
	// Decide lets the advisory findings call the decision model for
	// concepts without a stored answer.
	Decide bool
	Flags  session.Flags
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
	command.Flags().BoolVar(&options.Decide, "decide", false, "Ask the decision model for advisory findings missing from receipts (requires [decisions])")
	command.Flags().Float64Var(&options.Flags.BudgetUSD, "budget", 0, "Spend ceiling in US$ for --decide (default [decisions].budget_usd)")
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

	advisory, s, err := okfAdvisoryOptions(cmd, cfg, options, topicInfo)
	if err != nil {
		return err
	}
	if s != nil {
		defer s.WriteSummary(cmd.ErrOrStderr())
	}
	issues, err := runOKFCheck(commandContext(cmd), topicInfo.RootPath, kokf.CheckOptions{
		Types:    cfg.OKF.Types,
		Strict:   options.Strict,
		Advisory: advisory,
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

// okfAdvisoryOptions configures the advisory findings. Without --decide they
// are read from the bundle's receipts only and no model is called; with
// --decide a decision session over the OKF topic asks the missing ones once
// (receipts cache them for every later check).
func okfAdvisoryOptions(cmd *cobra.Command, cfg kconfig.Config, options *okfCheckOptions, topicInfo models.TopicInfo) (*kokf.AdvisoryOptions, *session.Session, error) {
	advisory := &kokf.AdvisoryOptions{Options: kokf.TypeOptions(cfg.OKF.Types, cfg.OKF.TypeDescription)}
	if !options.Decide {
		settings, err := contract.LoadSettings(topicInfo.RootPath)
		if err != nil {
			return nil, nil, fmt.Errorf("okf check: %w", err)
		}
		advisory.Threshold = decisions.Thresholds(nil).Merge(cfg.Decisions.Thresholds).Merge(settings.Decisions.Thresholds).Get("okf_type")
		advisory.Topic.Exclude = settings.Decisions.Exclude
		return advisory, nil, nil
	}

	s, err := openSession(cmd, "kb okf check", topicInfo.Slug, options.Flags)
	if err != nil {
		return nil, nil, err
	}
	advisory.Decider = s.Engine
	advisory.Topic = s.Ref
	advisory.Threshold = s.Threshold("okf_type")
	advisory.Extras = s.ExtraBanks(decisions.PurposeOKFType)
	return advisory, s, nil
}
