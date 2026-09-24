package corpus

import (
	"path"
	"strings"
)

// MatchGlob reports whether the slash path p matches pattern. Matching is
// segment-aware: `*`, `?` and `[...]` match within one path segment (as in
// path.Match) and a `**` segment matches zero or more whole segments, so
// `raw/audit-*/**` matches every file below any `raw/audit-...` folder. Leading
// `./` and surrounding slashes are ignored on both sides; a malformed pattern
// matches nothing.
func MatchGlob(pattern, p string) bool {
	patternSegments := globSegments(pattern)
	pathSegments := globSegments(p)
	if len(patternSegments) == 0 {
		return len(pathSegments) == 0
	}

	return matchSegments(patternSegments, pathSegments)
}

// ValidGlob reports whether every segment of pattern is well formed.
func ValidGlob(pattern string) bool {
	for _, segment := range globSegments(pattern) {
		if segment == "**" {
			continue
		}
		if _, err := path.Match(segment, ""); err != nil {
			return false
		}
	}

	return true
}

func globSegments(value string) []string {
	cleaned := strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	for strings.HasPrefix(cleaned, "./") {
		cleaned = strings.TrimPrefix(cleaned, "./")
	}
	cleaned = strings.Trim(cleaned, "/")
	if cleaned == "" {
		return nil
	}

	return strings.Split(cleaned, "/")
}

func matchSegments(pattern, segments []string) bool {
	for len(pattern) > 0 {
		head := pattern[0]
		if head == "**" {
			rest := pattern[1:]
			for len(rest) > 0 && rest[0] == "**" {
				rest = rest[1:]
			}
			if len(rest) == 0 {
				return true
			}
			for index := range len(segments) + 1 {
				if matchSegments(rest, segments[index:]) {
					return true
				}
			}
			return false
		}
		if len(segments) == 0 {
			return false
		}
		matched, err := path.Match(head, segments[0])
		if err != nil || !matched {
			return false
		}
		pattern, segments = pattern[1:], segments[1:]
	}

	return len(segments) == 0
}
