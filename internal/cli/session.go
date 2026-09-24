package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/session"
)

var sessionGetwd = os.Getwd

// openSession opens a decision-backed session for one topic. Tests replace
// it to point the session at fake servers.
var openSession = func(cmd *cobra.Command, command, topicSlug string, flags session.Flags) (*session.Session, error) {
	// An explicit --budget (0 included) is honored; an invalid one fails
	// before anything else happens.
	flags.BudgetSet = flags.BudgetSet || budgetFlagSet(cmd)
	if _, err := flags.Budget(); err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}
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
// command. --budget is validated while the flags are parsed, so an invalid
// amount fails before the command does anything.
func bindDecisionFlags(command *cobra.Command, flags *session.Flags) {
	command.Flags().Var(&budgetValue{flags: flags}, "budget", "Spend ceiling in US$ for this run; 0 = cached answers only (default [decisions].budget_usd)")
	command.Flags().StringVar(&flags.Decisions, "decisions", "", "Decision mode override for this run: shadow|apply")
}

// budgetFlagSet reports whether the command was given --budget.
func budgetFlagSet(cmd *cobra.Command) bool {
	return cmd != nil && cmd.Flags().Lookup("budget") != nil && cmd.Flags().Changed("budget")
}

// budgetValue is the --budget flag: it records that the flag was given (an
// explicit 0 is a cache-only run, not "unset") and rejects negative, NaN and
// infinite amounts at parse time.
type budgetValue struct {
	flags *session.Flags
}

func (v *budgetValue) String() string {
	if v.flags == nil {
		return "0"
	}
	return strconv.FormatFloat(v.flags.BudgetUSD, 'g', -1, 64)
}

func (v *budgetValue) Set(raw string) error {
	amount, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return fmt.Errorf("invalid amount %q", raw)
	}
	next := *v.flags
	next.BudgetUSD, next.BudgetSet = amount, true
	if _, err := next.Budget(); err != nil {
		return err
	}
	*v.flags = next
	return nil
}

func (v *budgetValue) Type() string { return "float" }
