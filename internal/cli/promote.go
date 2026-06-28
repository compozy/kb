package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	kokf "github.com/compozy/kb/internal/okf"
	ktopic "github.com/compozy/kb/internal/topic"
)

var runPromote = kokf.Promote
var runPromoteTopicInfo = ktopic.Info
var promoteGetwd = os.Getwd

func newPromoteCommand() *cobra.Command {
	var targetTopic string
	var conceptType string
	var description string

	command := &cobra.Command{
		Use:   "promote <wiki-doc>",
		Short: "Promote a wiki document into an OKF topic",
		Args:  cobra.ExactArgs(1),
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

			result, err := runPromote(commandContext(cmd), kokf.PromoteInput{
				SourceDocPath: args[0],
				VaultPath:     vaultPath,
				TargetTopic:   target,
				Type:          conceptType,
				Description:   description,
				Types:         cfg.OKF.Types,
			})
			if err != nil {
				return err
			}
			return writeJSON(cmd, result)
		},
	}

	command.Flags().StringVar(&targetTopic, "to", "", "Target OKF topic slug")
	command.Flags().StringVar(&conceptType, "type", "", "OKF concept type")
	command.Flags().StringVar(&description, "description", "", "OKF concept description")
	_ = command.MarkFlagRequired("to")
	_ = command.MarkFlagRequired("type")
	return command
}
