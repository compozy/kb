package cli

import (
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/review"
)

func TestSelectReviewItems(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	store := review.Open(root, nil)
	items := []review.Item{
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/a.md", Question: "role", Probability: 0.9},
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/b.md", Question: "role", Probability: 0.6},
		{Queue: review.QueueLink, Purpose: "link", Subject: "raw/a.md", Target: "wiki/concepts/X.md", Question: "should_link", Probability: 0.8},
		{Queue: review.QueueLink, Purpose: "mention", Subject: "raw/b.md", Target: "wiki/concepts/X.md", Question: "mention_sense", Probability: 0.95},
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		item.ID = review.ItemID(item.Queue, item.Subject, item.Target, item.Question)
		if _, err := store.Add(item); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, item.ID)
	}
	if _, err := store.Resolve(ids[1], review.StatusRejected); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		sel     reviewSelection
		want    []string
		wantErr string
	}{
		{name: "Should select pending ids", sel: reviewSelection{IDs: []string{ids[0], ids[2]}}, want: []string{ids[0], ids[2]}},
		{name: "Should refuse a resolved id", sel: reviewSelection{IDs: []string{ids[1]}}, wantErr: "already rejected"},
		{name: "Should refuse an id outside the selected queue", sel: reviewSelection{IDs: []string{ids[2]}, Queue: review.QueueRemove}, wantErr: "not the selected one"},
		{name: "Should refuse unknown ids", sel: reviewSelection{IDs: []string{"rv-nope"}}, wantErr: "unknown review item"},
		{name: "Should refuse mixing ids and bulk selectors", sel: reviewSelection{IDs: []string{ids[0]}, TopicWide: true}, wantErr: "not both"},
		{name: "Should require a bulk selector", sel: reviewSelection{Queue: review.QueueLink}, wantErr: "--above <p> or --topic-wide"},
		{name: "Should require a queue or purpose in bulk", sel: reviewSelection{TopicWide: true}, wantErr: "needs --queue or --all-purpose"},
		{name: "Should refuse unknown queues", sel: reviewSelection{Queue: "mystery", TopicWide: true}, wantErr: "unknown queue"},
		{name: "Should select a queue above a probability", sel: reviewSelection{Queue: review.QueueLink, Above: 0.9, AboveSet: true}, want: []string{ids[3]}},
		{name: "Should filter by purpose", sel: reviewSelection{Purpose: "link", Above: 0.75, AboveSet: true}, want: []string{ids[2]}},
		{name: "Should select every pending item of a queue topic-wide", sel: reviewSelection{Queue: review.QueueLink, TopicWide: true}, want: []string{ids[3], ids[2]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := selectReviewItems(review.Open(root, nil), tt.sel)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			gotIDs := make([]string, 0, len(got))
			for _, item := range got {
				gotIDs = append(gotIDs, item.ID)
			}
			if !slices.Equal(gotIDs, tt.want) {
				t.Fatalf("selected = %v, want %v", gotIDs, tt.want)
			}
		})
	}
}
