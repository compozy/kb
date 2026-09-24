package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
)

func newIngestBookmarksCommand() *cobra.Command {
	var topic string
	var batch string

	command := &cobra.Command{
		Use:   "bookmarks <path>",
		Short: "Ingest a bookmark-cluster markdown file into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest bookmarks", topic)
			if err != nil {
				return err
			}

			sourcePath, err := resolveInputFile("ingest bookmarks", args[0])
			if err != nil {
				return err
			}

			result, err := runIngest(commandContext(cmd), kingest.Options{
				VaultPath:  target.VaultPath,
				Topic:      target.TopicInfo.Slug,
				SourceKind: models.SourceKindBookmarkCluster,
				SourcePath: sourcePath,
				Registry:   newIngestRegistry(),
				Batch:      resolveIngestBatch("bookmarks", batch),
				// The bookmark file name is the cluster's label; its title is
				// only known after conversion, so it is not guessed here.
				Query: filepath.Base(sourcePath),
			})
			if err != nil {
				return fmt.Errorf("ingest bookmarks: %w", err)
			}

			return writeJSON(cmd, result)
		},
	}

	requireTopicFlag(command, &topic)
	addBatchFlag(command, &batch)

	return command
}
