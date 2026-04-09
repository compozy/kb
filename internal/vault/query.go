package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const topicMarkerFile = "CLAUDE.md"

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

// DiscoverVaultPath walks up from cwd until it finds a `.kodebase/vault` directory.
func DiscoverVaultPath(cwd string) (string, error) {
	resolvedCWD, err := resolveAbsolutePath(cwd)
	if err != nil {
		return "", fmt.Errorf("discover vault path: %w", err)
	}

	currentPath := resolvedCWD
	for {
		candidatePath := filepath.Join(currentPath, ".kodebase", "vault")
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
		"unable to find a vault from %s. walked up looking for .kodebase/vault/",
		resolvedCWD,
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

		topicPath := filepath.Join(vaultPath, topicSlug)
		if err := ensureDirectory(topicPath, "Topic path"); err != nil {
			return ResolvedVault{}, err
		}

		return ResolvedVault{
			VaultPath: vaultPath,
			TopicPath: topicPath,
			TopicSlug: topicSlug,
		}, nil
	}

	topics, err := listTopicDirectories(vaultPath)
	if err != nil {
		return ResolvedVault{}, fmt.Errorf("resolve vault query: %w", err)
	}

	switch len(topics) {
	case 0:
		return ResolvedVault{}, fmt.Errorf(
			"no topics were found in %s. expected child directories containing %s",
			vaultPath,
			topicMarkerFile,
		)
	case 1:
		return ResolvedVault{
			VaultPath: vaultPath,
			TopicPath: filepath.Join(vaultPath, topics[0]),
			TopicSlug: topics[0],
		}, nil
	default:
		return ResolvedVault{}, fmt.Errorf(
			"multiple topics were found in %s: %s. re-run with --topic <slug>",
			vaultPath,
			strings.Join(topics, ", "),
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

	topics, err := listTopicDirectories(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("list available topics: %w", err)
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

func listTopicDirectories(vaultPath string) ([]string, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("read vault path %q: %w", vaultPath, err)
	}

	topics := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		markerPath := filepath.Join(vaultPath, entry.Name(), topicMarkerFile)
		info, err := os.Stat(markerPath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("stat topic marker %q: %w", markerPath, err)
		}
		if info.IsDir() {
			continue
		}

		topics = append(topics, entry.Name())
	}

	sort.Strings(topics)
	return topics, nil
}

func isDirectoryPath(pathToCheck string) bool {
	info, err := os.Stat(pathToCheck)
	if err != nil {
		return false
	}

	return info.IsDir()
}
