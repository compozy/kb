package decisions

import (
	"regexp"
	"strings"
)

// Redacted replaces every secret removed by Redact.
const Redacted = "[REDACTED]"

// secretPatterns are obvious credential shapes (spec §14). They are
// best-effort: keeping credentials out of the state at the source is still
// the caller's job.
var secretPatterns = []*regexp.Regexp{
	// PEM private key blocks (RSA, EC, OPENSSH, ...).
	regexp.MustCompile(`-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z0-9 ]*PRIVATE KEY-----`),
	// OpenAI / OpenRouter / Anthropic style keys: sk-..., sk-or-v1-..., sk-proj-..., sk-ant-...
	regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{16,}`),
	// GitHub tokens.
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
	regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}`),
	// Slack tokens.
	regexp.MustCompile(`\bxox[abposr]-[A-Za-z0-9-]{10,}`),
	// AWS access key ids.
	regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`),
	// Google API keys.
	regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}`),
	// JSON web tokens.
	regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}`),
}

// bearerPattern keeps the scheme word and replaces the token.
var bearerPattern = regexp.MustCompile(`(?i)\b(bearer)\s+[A-Za-z0-9._~+/=-]{16,}`)

// assignmentPattern replaces the value of `api_key=...`, `password: ...`
// style assignments while keeping the name, so the text stays readable.
// Values shorter than 8 characters are kept (business text such as
// `token=invoice`).
var assignmentPattern = regexp.MustCompile(`(?i)\b(api[_-]?key|apikey|access[_-]?token|refresh[_-]?token|auth[_-]?token|secret[_-]?key|client[_-]?secret|private[_-]?key|password|passwd|secret|token)(\s*[:=]\s*)(["']?)[^\s"',;]{8,}`)

// secretFieldNames are object keys whose string values are always redacted.
var secretFieldNames = map[string]bool{
	"password": true, "passwd": true, "passphrase": true, "secret": true,
	"clientsecret": true, "token": true, "accesstoken": true, "refreshtoken": true,
	"apikey": true, "authorization": true, "privatekey": true, "credentials": true,
}

// RedactString replaces obvious secrets in s with Redacted.
func RedactString(s string) string {
	for _, pattern := range secretPatterns {
		s = pattern.ReplaceAllString(s, Redacted)
	}
	s = bearerPattern.ReplaceAllString(s, "${1} "+Redacted)
	return assignmentPattern.ReplaceAllString(s, "${1}${2}${3}"+Redacted)
}

// Redact returns a copy of v, as a generic JSON tree, with obvious secrets
// replaced by Redacted in every string value (API keys, bearer tokens, PEM
// private keys, `api_key=...` assignments) and in the string values of
// secret-looking keys (password, token, api_key, ...). A string in returns a
// string. A value that cannot be encoded as JSON is returned unchanged; the
// engine then rejects it as an invalid request.
func Redact(v any) any {
	if text, ok := v.(string); ok {
		return RedactString(text)
	}
	tree, err := toTree(v)
	if err != nil {
		return v
	}
	return redactTree(tree)
}

func redactTree(node any) any {
	switch value := node.(type) {
	case string:
		return RedactString(value)
	case map[string]any:
		for key, item := range value {
			if text, isString := item.(string); isString && text != "" && isSecretField(key) {
				value[key] = Redacted
				continue
			}
			value[key] = redactTree(item)
		}
		return value
	case []any:
		for index, item := range value {
			value[index] = redactTree(item)
		}
		return value
	default:
		return node
	}
}

func isSecretField(key string) bool {
	normalized := strings.NewReplacer("_", "", "-", "").Replace(strings.ToLower(key))
	return secretFieldNames[normalized]
}
