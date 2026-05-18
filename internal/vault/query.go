package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/topic"
)

// ResolvedVault identifies the vault root and topic directory selected for a command.
type ResolvedVault struct {
	VaultPath string `json:"vaultPath"`
	TopicPath string `json:"topicPath"`
	TopicSlug string `json:"topicSlug"`
}

// VaultQueryOptions control vault and topic discovery from CLI-style flags.
type VaultQueryOptions struct {
	Vault string
	Topic string
	CWD   string
}

// DiscoverVaultPath walks up from cwd until it finds kb.toml or a legacy
// `.kb/vault` directory.
func DiscoverVaultPath(cwd string) (string, error) {
	resolvedCWD, err := resolveAbsolutePath(cwd)
	if err != nil {
		return "", fmt.Errorf("discover vault path: %w", err)
	}

	if configPath, found, err := config.DiscoverProjectConfigPath(resolvedCWD); err != nil {
		return "", fmt.Errorf("discover vault path: %w", err)
	} else if found {
		cfg, err := config.Load(configPath)
		if err != nil {
			return "", fmt.Errorf("discover vault path: load project config %q: %w", configPath, err)
		}
		vaultRoot, err := config.ResolveVaultRoot(configPath, cfg.Vault)
		if err != nil {
			return "", fmt.Errorf("discover vault path: %w", err)
		}
		if isDirectoryPath(vaultRoot) {
			return vaultRoot, nil
		}
		return "", fmt.Errorf("configured vault root was not found or is not a directory: %s", vaultRoot)
	}

	currentPath := resolvedCWD
	for {
		candidatePath := filepath.Join(currentPath, ".kb", "vault")
		if isDirectoryPath(candidatePath) {
			return candidatePath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}

		currentPath = parentPath
	}

	return "", fmt.Errorf(
		"unable to find a vault from %s. walked up looking for %s or .kb/vault/",
		resolvedCWD,
		config.ProjectConfigFileName,
	)
}

// ResolveVaultQuery resolves the target vault and topic from the provided options.
func ResolveVaultQuery(options VaultQueryOptions) (ResolvedVault, error) {
	cwd, err := resolveAbsolutePath(options.CWD)
	if err != nil {
		return ResolvedVault{}, fmt.Errorf("resolve vault query: %w", err)
	}

	vaultPath, err := resolveVaultPath(options, cwd)
	if err != nil {
		return ResolvedVault{}, err
	}
	if err := ensureDirectory(vaultPath, "Vault path"); err != nil {
		return ResolvedVault{}, err
	}

	topicSlug := strings.TrimSpace(options.Topic)
	if options.Topic != "" {
		if topicSlug == "" {
			return ResolvedVault{}, fmt.Errorf("topic name is required when topic is specified")
		}

		topicInfo, err := topic.Resolve(vaultPath, topicSlug)
		if err != nil {
			return ResolvedVault{}, err
		}

		return ResolvedVault{
			VaultPath: vaultPath,
			TopicPath: topicInfo.RootPath,
			TopicSlug: topicInfo.Slug,
		}, nil
	}

	topics, err := topic.List(vaultPath)
	if err != nil {
		return ResolvedVault{}, fmt.Errorf("resolve vault query: %w", err)
	}

	switch len(topics) {
	case 0:
		return ResolvedVault{}, fmt.Errorf(
			"no topics were found in %s. expected child directories containing %s",
			vaultPath,
			"CLAUDE.md",
		)
	case 1:
		return ResolvedVault{
			VaultPath: vaultPath,
			TopicPath: topics[0].RootPath,
			TopicSlug: topics[0].Slug,
		}, nil
	default:
		topicSlugs := make([]string, 0, len(topics))
		for _, topicInfo := range topics {
			topicSlugs = append(topicSlugs, topicInfo.Slug)
		}
		return ResolvedVault{}, fmt.Errorf(
			"multiple topics were found in %s: %s. re-run with --topic <slug>",
			vaultPath,
			strings.Join(topicSlugs, ", "),
		)
	}
}

// ListAvailableTopics returns the marker-backed topic directories in deterministic order.
func ListAvailableTopics(options VaultQueryOptions) ([]string, error) {
	cwd, err := resolveAbsolutePath(options.CWD)
	if err != nil {
		return nil, fmt.Errorf("list available topics: %w", err)
	}

	vaultPath, err := resolveVaultPath(options, cwd)
	if err != nil {
		return nil, err
	}
	if err := ensureDirectory(vaultPath, "Vault path"); err != nil {
		return nil, err
	}

	topicInfos, err := topic.List(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("list available topics: %w", err)
	}

	topics := make([]string, 0, len(topicInfos))
	for _, topicInfo := range topicInfos {
		topics = append(topics, topicInfo.Slug)
	}
	return topics, nil
}

func resolveVaultPath(options VaultQueryOptions, cwd string) (string, error) {
	if strings.TrimSpace(options.Vault) != "" {
		absoluteVaultPath, err := filepath.Abs(options.Vault)
		if err != nil {
			return "", fmt.Errorf("resolve vault query: resolve vault path %q: %w", options.Vault, err)
		}
		return absoluteVaultPath, nil
	}

	vaultPath, err := DiscoverVaultPath(cwd)
	if err != nil {
		return "", fmt.Errorf("resolve vault query: %w", err)
	}

	return vaultPath, nil
}

func resolveAbsolutePath(pathValue string) (string, error) {
	if strings.TrimSpace(pathValue) == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("resolve cwd: %w", err)
		}
		pathValue = currentPath
	}

	absolutePath, err := filepath.Abs(pathValue)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path %q: %w", pathValue, err)
	}

	return absolutePath, nil
}

func ensureDirectory(pathToCheck, label string) error {
	if !isDirectoryPath(pathToCheck) {
		return fmt.Errorf("%s was not found or is not a directory: %s", label, pathToCheck)
	}
	return nil
}

func isDirectoryPath(pathToCheck string) bool {
	info, err := os.Stat(pathToCheck)
	if err != nil {
		return false
	}

	return info.IsDir()
}
