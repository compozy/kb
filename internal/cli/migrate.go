package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	kmigrate "github.com/compozy/kb/internal/migrate"
)

var runMigrateTranscripts = kmigrate.Transcripts

func newMigrateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Run explicit vault layout migrations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.AddCommand(newMigrateTranscriptsCommand())

	return command
}

func newMigrateTranscriptsCommand() *cobra.Command {
	var topic string

	command := &cobra.Command{
		Use:   "transcripts",
		Short: "Move legacy raw/transcripts markdown files into raw/youtube",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := resolveCommandVaultPath(cmd, topicGetwd, "migrate transcripts")
			if err != nil {
				return err
			}

			result, err := runMigrateTranscripts(vaultPath, topic, time.Now().UTC())
			if err != nil {
				return fmt.Errorf("migrate transcripts: %w", err)
			}

			return writeJSON(cmd, result)
		},
	}

	requireTopicFlag(command, &topic)

	return command
}
