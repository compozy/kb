package questions

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
)

//go:embed banks/*.json
var embeddedBanks embed.FS

// ExitOptions are the abstention options a choice question may offer. Each
// has a distinct meaning (see jev-engineering questions.md); every choice in
// a bank must carry at least one of them.
var ExitOptions = []string{"unknown", "none", "other", "not_applicable", "different"}

// ErrUnknownBank is returned by Load for an id with no embedded bank.
var ErrUnknownBank = errors.New("questions: unknown bank")

// Bank is a loaded question bank. Hash covers the whole bank file and enters
// every cache key; Version is recorded in receipts and the state record.
// A Bank is immutable after loading and safe for concurrent use.
type Bank struct {
	ID      string
	Version string
	Hash    string
	Purpose string
	Guard   string

	questions map[string]rawQuestion
	order     []string
}

// rawQuestion is one editorial bank entry. Instructions already carry the
// guard once the bank is loaded; Text keeps the unguarded instructions for
// checks that read the question wording.
type rawQuestion struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Instructions  any    `json:"instructions"`
	Criteria      any    `json:"criteria,omitempty"`
	Source        string `json:"source"`
	Applicability string `json:"applicability,omitempty"`
	// InjectedOptions marks a choice whose real options are supplied by the
	// caller through WithCriteria; the bank then only carries exit options.
	InjectedOptions bool `json:"injected_options,omitempty"`

	text any
}

type bankFile struct {
	ID        string        `json:"id"`
	Version   string        `json:"version"`
	Purpose   string        `json:"purpose"`
	Guard     string        `json:"guard"`
	Questions []rawQuestion `json:"questions"`
}

var (
	idPattern          = regexp.MustCompile(`^[a-z][a-z0-9_]*(\{(id|n)\}[a-z0-9_]*)*$`)
	placeholderPattern = regexp.MustCompile(`\{([a-z]+)\}`)

	cacheMu sync.Mutex
	cache   = map[string]*Bank{}
)

// Load returns the embedded bank with the given id. Banks are parsed and
// validated once and cached for the life of the process.
func Load(id string) (*Bank, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if bank, ok := cache[id]; ok {
		return bank, nil
	}
	if !idPattern.MatchString(id) || strings.Contains(id, "{") {
		return nil, fmt.Errorf("%w: %q", ErrUnknownBank, id)
	}
	data, err := embeddedBanks.ReadFile(path.Join("banks", id+".json"))
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrUnknownBank, id)
	}
	bank, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("load bank %q: %w", id, err)
	}
	if bank.ID != id {
		return nil, fmt.Errorf("load bank %q: file declares id %q", id, bank.ID)
	}
	cache[id] = bank
	return bank, nil
}

// MustLoad is Load for embedded banks known at compile time; it panics on
// error, which a unit test over every embedded bank rules out.
func MustLoad(id string) *Bank {
	bank, err := Load(id)
	if err != nil {
		panic(err)
	}
	return bank
}

// BankIDs lists the ids of every embedded bank, sorted.
func BankIDs() []string {
	entries, err := embeddedBanks.ReadDir("banks")
	if err != nil {
		return nil
	}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if name, ok := strings.CutSuffix(entry.Name(), ".json"); ok {
			ids = append(ids, name)
		}
	}
	sort.Strings(ids)
	return ids
}

// Parse validates one bank file and returns the loaded bank with the guard
// applied to every question and Hash set to the sha256 of the canonical bank
// JSON.
func Parse(data []byte) (*Bank, error) {
	var file bankFile
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&file); err != nil {
		return nil, fmt.Errorf("decode bank: %w", err)
	}
	if err := validateBankHeader(file); err != nil {
		return nil, err
	}
	hash, err := canonicalHash(data)
	if err != nil {
		return nil, err
	}
	bank := &Bank{
		ID:        file.ID,
		Version:   file.Version,
		Hash:      hash,
		Purpose:   file.Purpose,
		Guard:     file.Guard,
		questions: make(map[string]rawQuestion, len(file.Questions)),
		order:     make([]string, 0, len(file.Questions)),
	}
	for _, question := range file.Questions {
		if err := validateQuestion(question); err != nil {
			return nil, fmt.Errorf("bank %q: %w", file.ID, err)
		}
		if _, dup := bank.questions[question.ID]; dup {
			return nil, fmt.Errorf("bank %q: duplicate question id %q", file.ID, question.ID)
		}
		question.text = question.Instructions
		question.Instructions = applyGuard(file.Guard, question.Instructions)
		bank.questions[question.ID] = question
		bank.order = append(bank.order, question.ID)
	}
	return bank, nil
}

func validateBankHeader(file bankFile) error {
	switch {
	case strings.TrimSpace(file.ID) == "":
		return errors.New("bank id is required")
	case !idPattern.MatchString(file.ID) || strings.Contains(file.ID, "{"):
		return fmt.Errorf("bank id %q must be lowercase snake_case", file.ID)
	case strings.TrimSpace(file.Version) == "":
		return fmt.Errorf("bank %q: version is required", file.ID)
	case strings.TrimSpace(file.Purpose) == "":
		return fmt.Errorf("bank %q: purpose is required", file.ID)
	case strings.TrimSpace(file.Guard) == "":
		return fmt.Errorf("bank %q: guard is required", file.ID)
	case len(file.Questions) == 0:
		return fmt.Errorf("bank %q: at least one question is required", file.ID)
	}
	return nil
}

func validateQuestion(q rawQuestion) error {
	if strings.TrimSpace(q.ID) == "" {
		return errors.New("question id is required")
	}
	if !idPattern.MatchString(q.ID) {
		return fmt.Errorf("question id %q must be snake_case with only {id} or {n} placeholders", q.ID)
	}
	if strings.TrimSpace(q.Source) == "" {
		return fmt.Errorf("question %q: source is required", q.ID)
	}
	if err := validateInstructions(q.Instructions); err != nil {
		return fmt.Errorf("question %q: %w", q.ID, err)
	}
	if q.InjectedOptions && q.Type != TypeChoice {
		return fmt.Errorf("question %q: injected_options is only valid for choice", q.ID)
	}
	switch q.Type {
	case TypeNoul:
		return validateNoulCriteria(q)
	case TypeChoice:
		return validateChoiceCriteria(q)
	case TypeScore:
		return validateScoreCriteria(q)
	default:
		return fmt.Errorf("question %q: type must be noul, choice or score: %q", q.ID, q.Type)
	}
}

func validateInstructions(instructions any) error {
	switch value := instructions.(type) {
	case string:
		if strings.TrimSpace(value) == "" {
			return errors.New("instructions are empty")
		}
	case map[string]any:
		if len(value) == 0 {
			return errors.New("instructions object is empty")
		}
		if _, ok := value["guard"]; ok {
			return errors.New(`instructions object must not set "guard"; the loader adds it`)
		}
	case []any:
		if len(value) == 0 {
			return errors.New("instructions array is empty")
		}
	default:
		return errors.New("instructions must be a string, object or array")
	}
	return nil
}

func validateNoulCriteria(q rawQuestion) error {
	if q.Criteria == nil {
		return nil
	}
	criteria, ok := q.Criteria.(map[string]any)
	if !ok {
		return fmt.Errorf("question %q: noul criteria must be an object with true/false keys", q.ID)
	}
	for key := range criteria {
		if key != "true" && key != "false" {
			return fmt.Errorf("question %q: noul criteria keys must be true or false, got %q", q.ID, key)
		}
	}
	return nil
}

func validateChoiceCriteria(q rawQuestion) error {
	criteria, ok := q.Criteria.(map[string]any)
	if !ok {
		return fmt.Errorf("question %q: choice criteria must be an object of options", q.ID)
	}
	minOptions := 2
	if q.InjectedOptions {
		minOptions = 1
	}
	if len(criteria) < minOptions || len(criteria) > 255 {
		return fmt.Errorf("question %q: choice needs %d to 255 options, got %d", q.ID, minOptions, len(criteria))
	}
	hasExit := false
	for option := range criteria {
		if strings.TrimSpace(option) == "" {
			return fmt.Errorf("question %q: choice option names must be non-empty", q.ID)
		}
		if slices.Contains(ExitOptions, option) {
			hasExit = true
		}
	}
	if !hasExit {
		return fmt.Errorf("question %q: choice needs an exit option (%s)", q.ID, strings.Join(ExitOptions, ", "))
	}
	return nil
}

func validateScoreCriteria(q rawQuestion) error {
	levels, ok := q.Criteria.([]any)
	if !ok {
		return fmt.Errorf("question %q: score criteria must be an ordered array of levels", q.ID)
	}
	if len(levels) < 2 || len(levels) > 10 {
		return fmt.Errorf("question %q: score needs 2 to 10 levels, got %d", q.ID, len(levels))
	}
	for index, level := range levels {
		if text, isString := level.(string); isString && strings.TrimSpace(text) == "" {
			return fmt.Errorf("question %q: score level %d is empty", q.ID, index)
		}
		if level == nil {
			return fmt.Errorf("question %q: score level %d is null", q.ID, index)
		}
	}
	return nil
}

// applyGuard prefixes guard to instructions without destroying their shape:
// string → "guard text", object → {"guard": guard, ...}, array → [guard, ...].
func applyGuard(guard string, instructions any) any {
	switch value := instructions.(type) {
	case string:
		return guard + " " + value
	case map[string]any:
		guarded := make(map[string]any, len(value)+1)
		maps.Copy(guarded, value)
		guarded["guard"] = guard
		return guarded
	case []any:
		return append([]any{guard}, value...)
	default:
		return instructions
	}
}

// IDs returns the question ids (templates included) in bank order.
func (b *Bank) IDs() []string {
	return slices.Clone(b.order)
}

// Question returns the API-ready question for templateID with every {key}
// placeholder substituted from vars, in the id and in every string inside the
// instructions and criteria. The returned Q is a deep copy. An id placeholder
// without a value is an error.
func (b *Bank) Question(templateID string, vars map[string]string) (Q, error) {
	raw, ok := b.questions[templateID]
	if !ok {
		return Q{}, fmt.Errorf("bank %q: unknown question %q", b.ID, templateID)
	}
	for _, value := range vars {
		if strings.ContainsAny(value, "{}") {
			return Q{}, fmt.Errorf("bank %q: question %q: placeholder values must not contain braces: %q", b.ID, templateID, value)
		}
	}
	id := substitute(templateID, vars)
	if missing := placeholderPattern.FindString(id); missing != "" {
		return Q{}, fmt.Errorf("bank %q: question %q needs a value for %s", b.ID, templateID, missing)
	}
	return Q{
		ID:           id,
		Type:         raw.Type,
		Instructions: substituteAny(raw.Instructions, vars),
		Criteria:     substituteAny(raw.Criteria, vars),
	}, nil
}

// MustQuestion is Question for fixed ids and vars; it panics on error.
func (b *Bank) MustQuestion(templateID string, vars map[string]string) Q {
	q, err := b.Question(templateID, vars)
	if err != nil {
		panic(err)
	}
	return q
}

// WithCriteria returns a copy of q whose criteria are replaced by criteria,
// for choices whose options are built by the caller (concepts, OKF types).
// When both the bank criteria and the new criteria are option maps, the bank's
// exit options (see ExitOptions) are merged in unless the caller already set
// them, so an injected choice always keeps its way out.
func (b *Bank) WithCriteria(q Q, criteria any) Q {
	out := Q{ID: q.ID, Type: q.Type, Instructions: deepCopy(q.Instructions), Criteria: deepCopy(normalizeCriteria(criteria))}
	next, isMap := out.Criteria.(map[string]any)
	if !isMap || q.Type != TypeChoice {
		return out
	}
	current, isCurrentMap := q.Criteria.(map[string]any)
	if !isCurrentMap {
		return out
	}
	for option, description := range current {
		if _, set := next[option]; !set && slices.Contains(ExitOptions, option) {
			next[option] = deepCopy(description)
		}
	}
	return out
}

// normalizeCriteria turns typed option maps and level slices into the generic
// JSON shapes used everywhere else.
func normalizeCriteria(criteria any) any {
	switch value := criteria.(type) {
	case map[string]string:
		out := make(map[string]any, len(value))
		for key, item := range value {
			out[key] = item
		}
		return out
	case []string:
		out := make([]any, len(value))
		for index, item := range value {
			out[index] = item
		}
		return out
	default:
		return criteria
	}
}

func substitute(text string, vars map[string]string) string {
	if len(vars) == 0 || !strings.Contains(text, "{") {
		return text
	}
	return placeholderPattern.ReplaceAllStringFunc(text, func(match string) string {
		if value, ok := vars[match[1:len(match)-1]]; ok {
			return value
		}
		return match
	})
}

func substituteAny(value any, vars map[string]string) any {
	switch typed := value.(type) {
	case string:
		return substitute(typed, vars)
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, item := range typed {
			out[substitute(key, vars)] = substituteAny(item, vars)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for index, item := range typed {
			out[index] = substituteAny(item, vars)
		}
		return out
	default:
		return value
	}
}

func deepCopy(value any) any {
	return substituteAny(value, nil)
}

func canonicalHash(data []byte) (string, error) {
	var tree any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&tree); err != nil {
		return "", fmt.Errorf("canonicalize bank: %w", err)
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(tree); err != nil {
		return "", fmt.Errorf("canonicalize bank: %w", err)
	}
	sum := sha256.Sum256(bytes.TrimRight(buffer.Bytes(), "\n"))
	return hex.EncodeToString(sum[:]), nil
}

// LoadTopicExtras reads the topic's extra banks from
// <topicRoot>/.decisions/banks/*.json (same schema as the embedded banks;
// purpose required). A topic may add questions, never replace built-in ones:
// an extra question id that collides with a built-in question id of the same
// purpose is an error. A missing directory yields no banks.
func LoadTopicExtras(topicRoot string) ([]*Bank, error) {
	dir := filepath.Join(topicRoot, ".decisions", "banks")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read topic banks %q: %w", dir, err)
	}
	builtin, err := builtinIDsByPurpose()
	if err != nil {
		return nil, err
	}
	var banks []*Bank
	seen := map[string]string{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		file := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read topic bank %q: %w", file, err)
		}
		bank, err := Parse(data)
		if err != nil {
			return nil, fmt.Errorf("topic bank %q: %w", file, err)
		}
		for _, id := range bank.order {
			if slices.Contains(builtin[bank.Purpose], id) {
				return nil, fmt.Errorf("topic bank %q: question %q collides with a built-in %s question", file, id, bank.Purpose)
			}
			key := bank.Purpose + "\x00" + id
			if other, dup := seen[key]; dup {
				return nil, fmt.Errorf("topic bank %q: question %q is also defined in %q", file, id, other)
			}
			seen[key] = file
		}
		banks = append(banks, bank)
	}
	return banks, nil
}

func builtinIDsByPurpose() (map[string][]string, error) {
	out := map[string][]string{}
	for _, id := range BankIDs() {
		bank, err := Load(id)
		if err != nil {
			return nil, err
		}
		out[bank.Purpose] = append(out[bank.Purpose], bank.order...)
	}
	return out, nil
}
