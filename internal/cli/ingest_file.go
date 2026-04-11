package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/go-devstack/internal/convert"
	kingest "github.com/user/go-devstack/internal/ingest"
	"github.com/user/go-devstack/internal/models"
)

var newIngestRegistry = func() kingest.Registry {
	return convert.NewRegistry()
}

func newIngestFileCommand() *cobra.Command {
	var topic string

	command := &cobra.Command{
		Use:   "file <path>",
		Short: "Convert a local file and ingest it into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest file", topic)
			if err != nil {
				return err
			}

			sourcePath, err := resolveInputFile("ingest file", args[0])
			if err != nil {
				return err
			}

			result, err := runIngest(commandContext(cmd), kingest.Options{
				VaultPath:  target.VaultPath,
				Topic:      target.TopicInfo.Slug,
				SourceKind: models.SourceKindDocument,
				SourcePath: sourcePath,
				Registry:   newIngestRegistry(),
			})
			if err != nil {
				return fmt.Errorf("ingest file: %w", err)
			}

			return writeJSON(cmd, result)
		},
	}

	requireTopicFlag(command, &topic)

	return command
}
