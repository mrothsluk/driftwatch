package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

type webhookNotifier struct {
	url    string
	client *http.Client
}

func newWebhookNotifier(url string) *webhookNotifier {
	return &webhookNotifier{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type webhookPayload struct {
	HasDrift  bool              `json:"has_drift"`
	DriftCount int             `json:"drift_count"`
	Resources []string          `json:"resources,omitempty"`
}

// Notify sends a JSON payload to the configured webhook URL.
func (w *webhookNotifier) Notify(report *drift.Report) error {
	resources := make([]string, 0, len(report.Drifts))
	for _, d := range report.Drifts {
		resources = append(resources, d.ResourceID)
	}

	p := webhookPayload{
		HasDrift:   report.HasDrift(),
		DriftCount: len(report.Drifts),
		Resources:  resources,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook notify: marshal: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook notify: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
