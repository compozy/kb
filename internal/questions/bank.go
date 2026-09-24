package questions

// Bank is a loaded question bank. Hash covers the whole bank file and enters
// every cache key; Version is recorded in receipts and the state record.
// W1-A implements Load/Question/WithCriteria and the embedded banks.
type Bank struct {
	ID      string
	Version string
	Hash    string
	Purpose string
	Guard   string

	questions map[string]rawQuestion
	order     []string
}

type rawQuestion struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Instructions  any    `json:"instructions"`
	Criteria      any    `json:"criteria,omitempty"`
	Source        string `json:"source"`
	Applicability string `json:"applicability,omitempty"`
}
