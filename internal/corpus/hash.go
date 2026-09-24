package corpus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
)

// BodyHash is the sha256 hex digest of body after normalizing `\r\n` to `\n`.
func BodyHash(body string) string {
	sum := sha256.Sum256([]byte(strings.ReplaceAll(body, "\r\n", "\n")))
	return hex.EncodeToString(sum[:])
}

// ValueHash is the sha256 hex digest of the canonical JSON of the normalized
// value: strings are trimmed, lists keep their order ([]any of strings hashes
// like []string), numbers encode as JSON numbers (int 2 and float 2.0 hash
// alike), maps encode with sorted keys and times as their frontmatter text
// (`2006-01-02` at midnight, RFC 3339 otherwise). A value written by kb and
// the same value parsed back from YAML therefore hash identically.
func ValueHash(v any) string {
	encoded, err := json.Marshal(normalizeHashValue(v))
	if err != nil {
		encoded = fmt.Appendf(nil, "%#v", v)
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

func normalizeHashValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		return strings.TrimSpace(typed)
	case bool:
		return typed
	case []string:
		items := make([]string, len(typed))
		for index, item := range typed {
			items[index] = strings.TrimSpace(item)
		}
		return items
	case map[string]any:
		normalized := make(map[string]any, len(typed))
		for key, item := range typed {
			normalized[key] = normalizeHashValue(item)
		}
		return normalized
	case time.Time:
		if typed.Hour() == 0 && typed.Minute() == 0 && typed.Second() == 0 && typed.Nanosecond() == 0 {
			return typed.Format(frontmatter.DateLayout)
		}
		return typed.Format(time.RFC3339Nano)
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return json.Number(strconv.FormatInt(reflected.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return json.Number(strconv.FormatUint(reflected.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		return json.Number(strconv.FormatFloat(reflected.Float(), 'f', -1, 64))
	case reflect.Slice, reflect.Array:
		items := make([]any, reflected.Len())
		for index := range reflected.Len() {
			items[index] = normalizeHashValue(reflected.Index(index).Interface())
		}
		return items
	default:
		return value
	}
}
