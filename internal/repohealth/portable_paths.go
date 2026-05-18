package repohealth

import (
	"fmt"
	"strings"
)

var windowsReservedNames = map[string]struct{}{
	"CON":  {},
	"PRN":  {},
	"AUX":  {},
	"NUL":  {},
	"COM1": {},
	"COM2": {},
	"COM3": {},
	"COM4": {},
	"COM5": {},
	"COM6": {},
	"COM7": {},
	"COM8": {},
	"COM9": {},
	"LPT1": {},
	"LPT2": {},
	"LPT3": {},
	"LPT4": {},
	"LPT5": {},
	"LPT6": {},
	"LPT7": {},
	"LPT8": {},
	"LPT9": {},
}

const windowsInvalidChars = `<>:"\\|?*`

type PathViolation struct {
	Path   string
	Reason string
}

func (v PathViolation) String() string {
	return fmt.Sprintf("%s: %s", v.Path, v.Reason)
}

func FindWindowsPathViolations(paths []string) []PathViolation {
	var violations []PathViolation
	for _, candidate := range paths {
		if reason, ok := windowsPathViolation(candidate); ok {
			violations = append(violations, PathViolation{
				Path:   candidate,
				Reason: reason,
			})
		}
	}
	return violations
}

func FormatPathViolations(violations []PathViolation) string {
	if len(violations) == 0 {
		return ""
	}

	lines := make([]string, 0, len(violations))
	for _, violation := range violations {
		lines = append(lines, violation.String())
	}
	return strings.Join(lines, "\n")
}

func windowsPathViolation(path string) (string, bool) {
	if path == "" {
		return "path is empty", true
	}

	components := strings.Split(path, "/")
	for _, component := range components {
		if component == "" {
			return "path contains an empty component", true
		}
		if reason, ok := windowsComponentViolation(component); ok {
			return reason, true
		}
	}
	return "", false
}

func windowsComponentViolation(component string) (string, bool) {
	for _, r := range component {
		if r < 32 {
			return "component contains an ASCII control character", true
		}
		if strings.ContainsRune(windowsInvalidChars, r) {
			return fmt.Sprintf("component contains Windows-invalid character %q", r), true
		}
	}

	if strings.HasSuffix(component, " ") {
		return "component ends with a space", true
	}
	if strings.HasSuffix(component, ".") {
		return "component ends with a period", true
	}

	baseName := component
	if dot := strings.IndexRune(component, '.'); dot >= 0 {
		baseName = component[:dot]
	}
	if _, ok := windowsReservedNames[strings.ToUpper(baseName)]; ok {
		return fmt.Sprintf("component uses Windows-reserved device name %q", baseName), true
	}

	return "", false
}
