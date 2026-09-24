package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	kokf "github.com/compozy/kb/internal/okf"
	"github.com/compozy/kb/internal/session"
	ktopic "github.com/compozy/kb/internal/topic"
)

var runPromote = kokf.Promote
var runPromoteTopicInfo = ktopic.Info
var promoteGetwd = os.Getwd

func newPromoteCommand() *cobra.Command {
	var targetTopic string
	var conceptType string
	var description string
	var flags session.Flags

	command := &cobra.Command{
		Use:   "promote <wiki-doc>",
		Short: "Promote a wiki document into an OKF topic",
		Long: "Promote a compiled wiki document into an OKF topic. With an [okf].types vocabulary, " +
			"the decision model suggests the type when --type is omitted and warns when --type " +
			"disagrees with a confident suggestion; promote then requires [decisions].",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := resolveCommandVaultPath(cmd, promoteGetwd, "promote")
			if err != nil {
				return err
			}
			cfg, err := loadCLIConfig()
			if err != nil {
				return fmt.Errorf("promote: %w", err)
			}
			target, err := runPromoteTopicInfo(vaultPath, targetTopic)
			if err != nil {
				return fmt.Errorf("promote: %w", err)
			}

			input := kokf.PromoteInput{
				SourceDocPath: args[0],
				VaultPath:     vaultPath,
				TargetTopic:   target,
				Type:          conceptType,
				Description:   description,
				Types:         cfg.OKF.Types,
			}
			// A type choice needs options: without a vocabulary there is no
			// question to ask, so promote runs without the decision model
			// and --type stays required (checked by okf.Promote).
			if len(cfg.OKF.Types) > 0 {
				sourceTopic, err := kokf.SourceTopic(vaultPath, args[0])
				if err != nil {
					return fmt.Errorf("promote: %w", err)
				}
				s, err := openSession(cmd, "kb promote", sourceTopic, flags)
				if err != nil {
					return err
				}
				defer s.WriteSummary(cmd.ErrOrStderr())
				input.Suggester = kokf.EngineTypeSuggester{
					Decider: s.Engine,
					Topic:   s.Ref,
					Options: kokf.TypeOptions(cfg.OKF.Types, cfg.OKF.TypeDescription),
				}
				input.TypeThreshold = s.Threshold("okf_type")
			}

			result, err := runPromote(commandContext(cmd), input)
			if err != nil {
				return err
			}
			if result.TypeWarning != "" {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: %s\n", result.TypeWarning)
			}
			return writeJSON(cmd, result)
		},
	}

	command.Flags().StringVar(&targetTopic, "to", "", "Target OKF topic slug")
	command.Flags().StringVar(&conceptType, "type", "", "OKF concept type (suggested by the decision model when omitted)")
	command.Flags().StringVar(&description, "description", "", "OKF concept description")
	command.Flags().Float64Var(&flags.BudgetUSD, "budget", 0, "Spend ceiling in US$ for this run (default [decisions].budget_usd)")
	_ = command.MarkFlagRequired("to")
	return command
}
