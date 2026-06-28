package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/output"
	ktopic "github.com/compozy/kb/internal/topic"
)

var runTopicNew = ktopic.New
var runTopicNewWithMode = ktopic.NewWithMode
var runTopicList = ktopic.List
var runTopicInfo = ktopic.Info
var topicGetwd = os.Getwd

func newTopicCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "topic",
		Short: "Scaffold and inspect knowledge base topics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.AddCommand(
		newTopicNewCommand(),
		newTopicListCommand(),
		newTopicInfoCommand(),
	)

	return command
}

func newTopicNewCommand() *cobra.Command {
	var mode string
	command := &cobra.Command{
		Use:   "new <slug> <title> <domain>",
		Short: "Scaffold a new knowledge base topic",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := resolveCommandVaultPath(cmd, topicGetwd, "topic new")
			if err != nil {
				return err
			}

			topicMode, err := parseTopicMode(mode)
			if err != nil {
				return err
			}
			var info models.TopicInfo
			if topicMode == models.TopicModeWiki {
				info, err = runTopicNew(vaultPath, args[0], args[1], args[2])
			} else {
				info, err = runTopicNewWithMode(vaultPath, args[0], args[1], args[2], topicMode)
			}
			if err != nil {
				return err
			}

			return writeTopicInfoJSON(cmd, info)
		},
	}
	command.Flags().StringVar(&mode, "mode", string(models.TopicModeWiki), "Topic mode (wiki|okf)")
	return command
}

func newTopicListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List scaffolded knowledge base topics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := resolveCommandVaultPath(cmd, topicGetwd, "topic list")
			if err != nil {
				return err
			}

			topics, err := runTopicList(vaultPath)
			if err != nil {
				return err
			}

			rows := make([]map[string]any, 0, len(topics))
			for _, topic := range topics {
				rows = append(rows, map[string]any{
					"slug":     topic.Slug,
					"title":    topic.Title,
					"domain":   topic.Domain,
					"sources":  topic.SourceCount,
					"articles": topic.ArticleCount,
				})
			}

			_, err = cmd.OutOrStdout().Write([]byte(output.FormatOutput(output.FormatOptions{
				Format:  output.OutputFormatTable,
				Columns: []string{"slug", "title", "domain", "sources", "articles"},
				Data:    rows,
			})))
			if err != nil {
				return fmt.Errorf("topic list: write output: %w", err)
			}

			return nil
		},
	}
}

func parseTopicMode(value string) (models.TopicMode, error) {
	switch strings.TrimSpace(value) {
	case "", string(models.TopicModeWiki):
		return models.TopicModeWiki, nil
	case string(models.TopicModeOKF):
		return models.TopicModeOKF, nil
	default:
		return "", fmt.Errorf(`invalid --mode %q: expected "wiki" or "okf"`, value)
	}
}

func newTopicInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info <slug>",
		Short: "Show metadata for one knowledge base topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultPath, err := resolveCommandVaultPath(cmd, topicGetwd, "topic info")
			if err != nil {
				return err
			}

			info, err := runTopicInfo(vaultPath, args[0])
			if err != nil {
				return err
			}

			return writeTopicInfoJSON(cmd, info)
		},
	}
}

func writeTopicInfoJSON(cmd *cobra.Command, info models.TopicInfo) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(info)
}
