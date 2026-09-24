package contract

import (
	"fmt"
	"path"
	"strings"
)

// ValidatePattern reports whether a topic-relative glob compiles: it must be
// non-empty, relative, free of `.`/`..` segments, use `**` only as a whole
// segment, and every other segment must be a valid path.Match pattern.
func ValidatePattern(pattern string) error {
	cleaned := strings.TrimSpace(pattern)
	if cleaned == "" {
		return fmt.Errorf("empty glob")
	}
	if strings.HasPrefix(cleaned, "/") || strings.Contains(cleaned, `\`) {
		return fmt.Errorf("glob %q must be a topic-relative path with forward slashes", pattern)
	}
	for segment := range strings.SplitSeq(cleaned, "/") {
		switch {
		case segment == "":
			return fmt.Errorf("glob %q has an empty segment", pattern)
		case segment == "." || segment == "..":
			return fmt.Errorf("glob %q cannot contain %q segments", pattern, segment)
		case segment == "**":
			continue
		case strings.Contains(segment, "**"):
			return fmt.Errorf("glob %q uses ** inside a segment; ** must be a whole segment", pattern)
		}
		if _, err := path.Match(segment, ""); err != nil {
			return fmt.Errorf("glob %q: segment %q: %w", pattern, segment, err)
		}
	}
	return nil
}

// MatchPath matches a topic-relative path against a glob with segment
// semantics: `*` matches within one segment, `?` one rune, `[...]` a class,
// and a `**` segment any number of segments (including none). Invalid
// patterns never match.
func MatchPath(pattern, topicRelPath string) bool {
	if ValidatePattern(pattern) != nil {
		return false
	}
	target := strings.TrimPrefix(strings.ReplaceAll(strings.TrimSpace(topicRelPath), `\`, "/"), "./")
	if target == "" {
		return false
	}
	return matchSegments(strings.Split(strings.TrimSpace(pattern), "/"), strings.Split(target, "/"))
}

func matchSegments(pattern, target []string) bool {
	for len(pattern) > 0 {
		if pattern[0] == "**" {
			rest := pattern[1:]
			for skip := 0; skip <= len(target); skip++ {
				if matchSegments(rest, target[skip:]) {
					return true
				}
			}
			return false
		}
		if len(target) == 0 {
			return false
		}
		matched, err := path.Match(pattern[0], target[0])
		if err != nil || !matched {
			return false
		}
		pattern, target = pattern[1:], target[1:]
	}
	return len(target) == 0
}
