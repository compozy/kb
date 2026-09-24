package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/output"
	"github.com/compozy/kb/internal/session"
)

var (
	runDraftVocabulary  = classify.DraftVocabulary
	runAcceptVocabulary = classify.AcceptVocabulary
)

type topicVocabularyOptions struct {
	Draft  bool
	Accept bool
	Flags  session.Flags
}

func newTopicVocabularyCommand() *cobra.Command {
	options := &topicVocabularyOptions{}
	command := &cobra.Command{
		Use:   "vocabulary <topic> --draft|--accept",
		Short: "Draft or accept the concept vocabulary of a topic without articles",
		Long: "--draft asks the generation model for 12 to 40 concepts (title + criterion) from the selection " +
			"contract and the source summaries, writes them to topic.yaml under vocabulary_draft and prints them " +
			"with how many source summaries mention each. Edit the draft, then --accept creates one stub article " +
			"per concept under wiki/concepts/; the next `kb classify <topic> --only-missing` fills concepts.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTopicVocabularyCommand(cmd, options, args[0])
		},
	}
	command.Flags().BoolVar(&options.Draft, "draft", false, "Generate a vocabulary draft into topic.yaml")
	command.Flags().BoolVar(&options.Accept, "accept", false, "Create one stub article per vocabulary_draft item")
	bindDecisionFlags(command, &options.Flags)
	return command
}

func runTopicVocabularyCommand(cmd *cobra.Command, options *topicVocabularyOptions, topicSlug string) error {
	if options.Draft == options.Accept {
		return fmt.Errorf("topic vocabulary: pass exactly one of --draft or --accept")
	}
	s, err := openSession(cmd, "kb topic vocabulary", topicSlug, options.Flags)
	if err != nil {
		return err
	}
	stdout := cmd.OutOrStdout()

	if options.Accept {
		created, err := runAcceptVocabulary(s)
		if err != nil {
			return fmt.Errorf("topic vocabulary: %w", err)
		}
		for _, rel := range created {
			if _, err := fmt.Fprintln(stdout, rel); err != nil {
				return fmt.Errorf("topic vocabulary: write output: %w", err)
			}
		}
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "created %d stub articles; next: kb classify %s --only-missing\n", len(created), s.Topic.Slug)
		return nil
	}

	items, err := runDraftVocabulary(cmd.Context(), s)
	s.WriteSummary(cmd.ErrOrStderr())
	if err != nil {
		return fmt.Errorf("topic vocabulary: %w", err)
	}
	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rows = append(rows, map[string]any{"title": item.Title, "mentions": item.Mentions, "criterion": item.Criterion})
	}
	if _, err := io.WriteString(stdout, output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"title", "mentions", "criterion"},
		Data:    rows,
	})); err != nil {
		return fmt.Errorf("topic vocabulary: write output: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "wrote %d concepts to topic.yaml vocabulary_draft; edit them, then run: kb topic vocabulary %s --accept\n", len(items), s.Topic.Slug)
	return nil
}
