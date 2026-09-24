package questions

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Compose returns one bank that asks the questions of every member bank in a
// single request (for example relevance + quality for the shared gate
// judgment). The composite's ID and Version are the members' joined by "+",
// its Hash is the sha256 of the members' hashes in order (so it changes when
// any member changes), and its Purpose and Guard are the first member's.
// Question and WithCriteria resolve a template id against the member that
// owns it; when two members define the same template id the first member
// wins. Nil members are ignored; Compose of no bank returns nil and Compose
// of one bank returns that bank unchanged.
func Compose(banks ...*Bank) *Bank {
	members := make([]*Bank, 0, len(banks))
	for _, bank := range banks {
		if bank != nil {
			members = append(members, bank)
		}
	}
	switch len(members) {
	case 0:
		return nil
	case 1:
		return members[0]
	}

	ids := make([]string, len(members))
	versions := make([]string, len(members))
	hasher := sha256.New()
	composite := &Bank{
		Purpose:   members[0].Purpose,
		Guard:     members[0].Guard,
		questions: map[string]rawQuestion{},
	}
	for index, member := range members {
		ids[index] = member.ID
		versions[index] = member.Version
		hasher.Write([]byte(member.Hash))
		hasher.Write([]byte{0})
		for _, id := range member.order {
			if _, taken := composite.questions[id]; taken {
				continue
			}
			composite.questions[id] = member.questions[id]
			composite.order = append(composite.order, id)
		}
	}
	composite.ID = strings.Join(ids, "+")
	composite.Version = strings.Join(versions, "+")
	composite.Hash = hex.EncodeToString(hasher.Sum(nil))
	return composite
}
