package review

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// ListView is the data behind `kb review list`.
type ListView struct {
	// Counts are the pending items per queue.
	Counts map[string]int `json:"counts"`
	// Items are the pending items shown.
	Items []ListItem `json:"items"`
}

// ListItem is one pending item with its subject title resolved.
type ListItem struct {
	Item
	Title string `json:"title,omitempty"`
}

// BuildListView collects the pending items of one queue ("" = all) with
// their subject titles (title returns "" when unknown).
func BuildListView(store *Store, queue string, title func(subject string) string) (ListView, error) {
	counts, err := store.PendingCounts()
	if err != nil {
		return ListView{}, err
	}
	items, err := store.Pending(queue)
	if err != nil {
		return ListView{}, err
	}
	view := ListView{Counts: counts, Items: make([]ListItem, 0, len(items))}
	for _, item := range items {
		entry := ListItem{Item: item}
		if title != nil {
			entry.Title = title(item.Subject)
		}
		view.Items = append(view.Items, entry)
	}
	return view, nil
}

// printer writes formatted text and keeps the first error.
type printer struct {
	w   io.Writer
	err error
}

func (p *printer) printf(format string, args ...any) {
	if p.err == nil {
		_, p.err = fmt.Fprintf(p.w, format, args...)
	}
}

// table writes one aligned table through a tabwriter.
func (p *printer) table(write func(t *printer)) {
	if p.err != nil {
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	inner := &printer{w: tw}
	write(inner)
	if inner.err != nil {
		p.err = inner.err
		return
	}
	p.err = tw.Flush()
}

// RenderList writes the queue counts and the pending items as tables.
func RenderList(w io.Writer, view ListView) error {
	p := &printer{w: w}
	p.table(func(t *printer) {
		t.printf("QUEUE\tPENDING\n")
		total := 0
		for _, queue := range Queues {
			t.printf("%s\t%d\n", queue, view.Counts[queue])
			total += view.Counts[queue]
		}
		t.printf("total\t%d\n", total)
	})
	if len(view.Items) == 0 {
		p.printf("\nNo pending items.\n")
		return wrapRender(p.err)
	}
	p.printf("\n")
	p.table(func(t *printer) {
		t.printf("ID\tQUEUE\tSUBJECT\tTARGET\tQUESTION\tP\tEVIDENCE\n")
		for _, item := range view.Items {
			subject := item.Subject
			if item.Title != "" {
				subject = item.Title + " (" + item.Subject + ")"
			}
			t.printf("%s\t%s\t%s\t%s\t%s\t%.2f\t%s\n",
				item.ID, item.Queue, cell(subject), cell(item.Target), cell(item.Question), item.Probability, cell(item.Evidence))
		}
	})
	return wrapRender(p.err)
}

// RenderCalibration writes a calibration report: one block per purpose with
// label counts, metrics at the current and chosen thresholds on dev and
// holdout, imported-label metrics and the dev sweep.
func RenderCalibration(w io.Writer, report Report) error {
	p := &printer{w: w}
	for index, purpose := range report.Purposes {
		if index > 0 {
			p.printf("\n")
		}
		p.printf("%s — %s\n", purpose.Purpose, purpose.Quantity)
		p.printf("labels: dev %d, holdout %d (imported %d, demoted to dev %d, without receipt %d)\n",
			purpose.DevLabels, purpose.HoldoutLabels, purpose.ImportedLabels, purpose.Demoted, purpose.Unjoined)
		if purpose.Recommended {
			p.printf("recommendation: %s = %.2f (current %.2f)\n", purpose.ApplyName, purpose.Chosen, purpose.Current)
		} else {
			p.printf("recommendation: none — %s\n", purpose.Note)
		}
		p.table(func(t *printer) {
			t.printf("SPLIT\tTHRESHOLD\tN\tPRECISION\tRECALL\tCOVERAGE\tREVIEW RATE\n")
			row := func(split string, threshold float64, m Metrics) {
				t.printf("%s\t%.2f\t%d\t%.2f\t%.2f\t%.2f\t%.2f\n", split, threshold, m.N, m.Precision, m.Recall, m.Coverage, m.ReviewRate)
			}
			row("dev (current)", purpose.Current, purpose.DevCurrent)
			row("holdout (current)", purpose.Current, purpose.HoldoutCurrent)
			chosen := purpose.Current
			if purpose.Recommended {
				chosen = purpose.Chosen
				row("dev (chosen)", chosen, purpose.DevChosen)
				row("holdout (chosen)", chosen, purpose.HoldoutChosen)
			}
			if purpose.ImportedLabels > 0 {
				row("imported dev", chosen, purpose.Imported.Dev)
				row("imported holdout", chosen, purpose.Imported.Holdout)
			}
		})
		p.table(func(t *printer) {
			t.printf("SWEEP (dev)\tPREDICTED\tPRECISION\tRECALL\n")
			for _, point := range purpose.Sweep {
				t.printf("%.2f\t%d\t%.2f\t%.2f\n", point.Threshold, point.Dev.Predicted, point.Dev.Precision, point.Dev.Recall)
			}
		})
	}
	return wrapRender(p.err)
}

func wrapRender(err error) error {
	if err != nil {
		return fmt.Errorf("review: render: %w", err)
	}
	return nil
}

func cell(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if value == "" {
		return "-"
	}
	if runes := []rune(value); len(runes) > 80 {
		return string(runes[:77]) + "..."
	}
	return value
}
