package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/kb/internal/vault"
)

const rootVaultFlagName = "vault"

func bindRootPersistentFlags(command *cobra.Command) {
	command.PersistentFlags().String(
		rootVaultFlagName,
		"",
		"Vault root path. Commands that read existing topics auto-discover .kb/vault/ from the current working directory when omitted.",
	)
}

func commandVaultValue(cmd *cobra.Command, current string) string {
	trimmedCurrent := strings.TrimSpace(current)
	if trimmedCurrent != "" {
		return trimmedCurrent
	}

	value, err := cmd.Flags().GetString(rootVaultFlagName)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(value)
}

func resolveCommandVaultPath(cmd *cobra.Command, getwd func() (string, error), action string) (string, error) {
	if explicitVaultPath := commandVaultValue(cmd, ""); explicitVaultPath != "" {
		absoluteVaultPath, err := filepath.Abs(explicitVaultPath)
		if err != nil {
			return "", fmt.Errorf("%s: resolve vault path %q: %w", action, explicitVaultPath, err)
		}
		return absoluteVaultPath, nil
	}

	cwd, err := getwd()
	if err != nil {
		return "", fmt.Errorf("%s: resolve cwd: %w", action, err)
	}

	discoveredVaultPath, err := vault.DiscoverVaultPath(cwd)
	if err != nil {
		return "", fmt.Errorf("%s: %w", action, err)
	}

	return discoveredVaultPath, nil
}
