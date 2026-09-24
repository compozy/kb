package decisions

import (
	"math"
	"sync"
)

// Budget is a per-run spend ceiling in US$, shared by the decision engine and
// the generation client. It is safe for concurrent use. A nil *Budget is
// unlimited.
type Budget struct {
	mu    sync.Mutex
	limit float64
	spent float64
}

// NewBudget returns a budget with the given ceiling. A ceiling of 0 allows no
// new calls, which makes a cache-only run.
func NewBudget(limitUSD float64) *Budget {
	if math.IsNaN(limitUSD) || limitUSD < 0 {
		limitUSD = 0
	}
	return &Budget{limit: limitUSD}
}

// Limit returns the ceiling.
func (b *Budget) Limit() float64 {
	if b == nil {
		return math.Inf(1)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.limit
}

// Spent returns what has been charged so far.
func (b *Budget) Spent() float64 {
	if b == nil {
		return 0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.spent
}

// Remaining returns the ceiling minus what was spent, never below zero.
func (b *Budget) Remaining() float64 {
	if b == nil {
		return math.Inf(1)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return math.Max(0, b.limit-b.spent)
}

// Exceeded reports whether no new call may start.
func (b *Budget) Exceeded() bool {
	if b == nil {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.spent >= b.limit
}

// Charge records a cost. Negative, NaN and infinite costs are ignored.
func (b *Budget) Charge(cost float64) {
	if b == nil || math.IsNaN(cost) || math.IsInf(cost, 0) || cost <= 0 {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.spent += cost
}
