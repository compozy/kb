package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/session"
)

var sessionGetwd = os.Getwd

// openSession opens a decision-backed session for one topic. Tests replace
// it to point the session at fake servers.
var openSession = func(cmd *cobra.Command, command, topicSlug string, flags session.Flags) (*session.Session, error) {
	vaultPath, err := resolveCommandVaultPath(cmd, sessionGetwd, command)
	if err != nil {
		return nil, err
	}
	cfg, err := loadCLIConfig()
	if err != nil {
		return nil, err
	}
	s, err := session.Open(session.Options{
		Config:    cfg,
		VaultPath: vaultPath,
		Topic:     topicSlug,
		Command:   command,
		Flags:     flags,
	})
	if err != nil {
		return nil, err
	}
	s.PrivacyNotice(cmd.ErrOrStderr())
	return s, nil
}

// bindDecisionFlags registers --budget and --decisions on a decision-backed
// command.
func bindDecisionFlags(command *cobra.Command, flags *session.Flags) {
	command.Flags().Float64Var(&flags.BudgetUSD, "budget", 0, "Spend ceiling in US$ for this run (default [decisions].budget_usd)")
	command.Flags().StringVar(&flags.Decisions, "decisions", "", "Decision mode override for this run: shadow|apply")
}
