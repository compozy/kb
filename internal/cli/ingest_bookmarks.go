package cli

import (
	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/models"
)

// newIngestBookmarksCommand ingests one bookmark-cluster markdown file as
// one source. The cluster is already collected (its URLs are not fetched),
// so there is no pre-fetch stage: the post-fetch gates judge the cluster
// document as a whole (dedupe by body hash, code checks without the thin
// rule, quality + relevance, near-duplicate). The file name is the
// ingest_query (the cluster label).
func newIngestBookmarksCommand() *cobra.Command {
	var topic string
	var batch string
	var flags ingestFlags

	command := &cobra.Command{
		Use:   "bookmarks <path>",
		Short: "Ingest a bookmark-cluster markdown file into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIngestLocalFile(cmd, "bookmarks", args[0], topic, batch, flags, models.SourceKindBookmarkCluster)
		},
	}

	requireTopicFlag(command, &topic)
	addBatchFlag(command, &batch)
	bindIngestFlags(command, &flags, false)

	return command
}
