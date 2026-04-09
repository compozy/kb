package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/go-devstack/internal/qmd"
	"github.com/user/go-devstack/internal/vault"
)

type indexCommandClient interface {
	Status(ctx context.Context) (qmd.IndexStatus, error)
	Index(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error)
}

type indexCommandOptions struct {
	Vault      string
	Topic      string
	Name       string
	Embed      bool
	Context    string
	ForceEmbed bool
}

var newIndexClient = func() indexCommandClient {
	return qmd.NewClient()
}

var resolveIndexVaultQuery = vault.ResolveVaultQuery
var indexGetwd = os.Getwd

func newIndexCommand() *cobra.Command {
	options := &indexCommandOptions{
		Embed: true,
	}

	command := &cobra.Command{
		Use:     "index",
		Aliases: []string{"index-vault"},
		Short:   "Create or update a QMD collection for a generated vault topic",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIndexCommand(cmd, options)
		},
	}

	flags := command.Flags()
	flags.StringVar(&options.Vault, "vault", "", "Vault root path")
	flags.StringVar(&options.Topic, "topic", "", "Topic slug inside the vault")
	flags.StringVar(&options.Name, "name", "", "Override the derived QMD collection name")
	flags.BoolVar(&options.Embed, "embed", true, "Run embedding after syncing files")
	flags.StringVar(&options.Context, "context", "", "Attach human-written collection context to improve search relevance")
	flags.BoolVar(&options.ForceEmbed, "force-embed", false, "Force re-embedding all documents")

	return command
}

func runIndexCommand(cmd *cobra.Command, options *indexCommandOptions) error {
	if options.ForceEmbed && !options.Embed {
		return fmt.Errorf("--force-embed cannot be used together with --embed=false")
	}

	cwd, err := indexGetwd()
	if err != nil {
		return fmt.Errorf("index: resolve cwd: %w", err)
	}

	resolvedVault, err := resolveIndexVaultQuery(vault.VaultQueryOptions{
		CWD:   cwd,
		Topic: strings.TrimSpace(options.Topic),
		Vault: strings.TrimSpace(options.Vault),
	})
	if err != nil {
		return fmt.Errorf("index: %w", err)
	}

	collectionName := strings.TrimSpace(options.Name)
	if collectionName == "" {
		collectionName = resolvedVault.TopicSlug
	}

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	client := newIndexClient()
	status, err := client.Status(ctx)
	if err != nil {
		return wrapQMDCommandError("index", err)
	}

	result, err := client.Index(ctx, qmd.IndexOptions{
		Operation:      chooseIndexOperation(status, collectionName),
		VaultPath:      resolvedVault.TopicPath,
		CollectionName: collectionName,
		Context:        strings.TrimSpace(options.Context),
		Embed:          options.Embed,
		ForceEmbed:     options.ForceEmbed,
	})
	if err != nil {
		return wrapQMDCommandError("index", err)
	}

	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	return encoder.Encode(indexResultPayload{
		CollectionName: collectionName,
		EmbedRequested: options.Embed,
		EmbedResult:    result.EmbedResult,
		ForceEmbed:     options.ForceEmbed,
		Status: indexStatusPayload{
			Collection:     findCollectionStatus(result.Status.Collections, collectionName),
			HasVectorIndex: result.Status.HasVectorIndex,
			NeedsEmbedding: result.Status.NeedsEmbedding,
			TotalDocuments: result.Status.TotalDocuments,
		},
		TopicPath:    resolvedVault.TopicPath,
		TopicSlug:    resolvedVault.TopicSlug,
		UpdateResult: result.UpdateResult,
		VaultPath:    resolvedVault.VaultPath,
	})
}

type indexResultPayload struct {
	CollectionName string             `json:"collectionName"`
	EmbedRequested bool               `json:"embedRequested"`
	EmbedResult    *qmd.EmbedResult   `json:"embedResult,omitempty"`
	ForceEmbed     bool               `json:"forceEmbed"`
	Status         indexStatusPayload `json:"status"`
	TopicPath      string             `json:"topicPath"`
	TopicSlug      string             `json:"topicSlug"`
	UpdateResult   qmd.UpdateResult   `json:"updateResult"`
	VaultPath      string             `json:"vaultPath"`
}

type indexStatusPayload struct {
	Collection     *qmd.CollectionInfo `json:"collection,omitempty"`
	HasVectorIndex bool                `json:"hasVectorIndex"`
	NeedsEmbedding int                 `json:"needsEmbedding"`
	TotalDocuments int                 `json:"totalDocuments"`
}

func chooseIndexOperation(status qmd.IndexStatus, collectionName string) qmd.IndexOperation {
	for _, collection := range status.Collections {
		if collection.Name == collectionName {
			return qmd.IndexOperationUpdate
		}
	}

	return qmd.IndexOperationAdd
}

func findCollectionStatus(collections []qmd.CollectionInfo, collectionName string) *qmd.CollectionInfo {
	for _, collection := range collections {
		if collection.Name != collectionName {
			continue
		}

		matched := collection
		return &matched
	}

	return nil
}
