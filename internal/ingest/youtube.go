package ingest

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/topic"
)

var youtubeVideoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// ExistingYouTubeVideoIDs returns the set of YouTube video IDs already ingested
// under the topic's raw/youtube directory. It powers resumable channel ingestion
// by reading each raw document's video_id frontmatter field and the id embedded in
// its source_url. A missing raw/youtube directory yields an empty set.
func ExistingYouTubeVideoIDs(vaultPath, topicRef string) (map[string]struct{}, error) {
	topicInfo, err := topic.Resolve(vaultPath, topicRef)
	if err != nil {
		return nil, fmt.Errorf("existing youtube ids: resolve topic: %w", err)
	}

	directory := filepath.Join(topicInfo.RootPath, "raw", "youtube")
	entries, err := os.ReadDir(directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]struct{}{}, nil
		}
		return nil, fmt.Errorf("existing youtube ids: read %q: %w", directory, err)
	}

	ids := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("existing youtube ids: read %q: %w", path, err)
		}
		values, _, err := frontmatter.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("existing youtube ids: parse %q: %w", path, err)
		}
		if id := strings.TrimSpace(frontmatter.GetString(values, "video_id")); youtubeVideoIDPattern.MatchString(id) {
			ids[id] = struct{}{}
		}
		if id := videoIDFromURL(frontmatter.GetString(values, "source_url")); id != "" {
			ids[id] = struct{}{}
		}
	}
	return ids, nil
}

func videoIDFromURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	if value := strings.TrimSpace(parsed.Query().Get("v")); youtubeVideoIDPattern.MatchString(value) {
		return value
	}
	segment := strings.Trim(parsed.Path, "/")
	if index := strings.LastIndex(segment, "/"); index >= 0 {
		segment = segment[index+1:]
	}
	if youtubeVideoIDPattern.MatchString(segment) {
		return segment
	}
	return ""
}
