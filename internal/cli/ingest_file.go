package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/convert"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
)

var newIngestRegistry = func() kingest.Registry {
	return convert.NewRegistry()
}

func newIngestFileCommand() *cobra.Command {
	var topic string
	var batch string
	var flags ingestFlags

	command := &cobra.Command{
		Use:   "file <path>",
		Short: "Convert a local file and ingest it into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIngestLocalFile(cmd, "file", args[0], topic, batch, flags, models.SourceKindDocument)
		},
	}

	requireTopicFlag(command, &topic)
	addBatchFlag(command, &batch)
	bindIngestFlags(command, &flags, false)

	return command
}

// runIngestLocalFile ingests one converted local file (file, bookmarks)
// through the gates. Local files have no URL, so stage 1 compares the body
// hash and the post-fetch stages 3, 5 and 6 apply; the file name is the
// ingest_query.
func runIngestLocalFile(cmd *cobra.Command, command, sourceArg, topicSlug, batch string, flags ingestFlags, kind models.SourceKind) error {
	action := "ingest " + command
	target, err := resolveIngestTarget(cmd, action, topicSlug)
	if err != nil {
		return err
	}
	sourcePath, err := resolveInputFile(action, sourceArg)
	if err != nil {
		return err
	}
	query := filepath.Base(sourcePath)
	runBatch := resolveIngestBatch(command, batch)
	runner, err := openIngestRunner(cmd, command, target.TopicInfo.Slug, flags, runBatch, query)
	if err != nil {
		return err
	}
	result, err := runner.Ingest(commandContext(cmd), kingest.Options{
		VaultPath:  target.VaultPath,
		Topic:      target.TopicInfo.Slug,
		SourceKind: kind,
		SourcePath: sourcePath,
		Registry:   newIngestRegistry(),
		Batch:      runBatch,
		Query:      query,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	if err := writeJSON(cmd, result); err != nil {
		return err
	}
	return finishIngestRun(cmd, runner)
}
