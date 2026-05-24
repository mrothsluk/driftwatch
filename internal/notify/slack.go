package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/driftwatch/internal/drift"
)

type slackNotifier struct {
	webhookURL string
	client     *http.Client
}

func newSlackNotifier(webhookURL string) *slackNotifier {
	return &slackNotifier{webhookURL: webhookURL, client: &http.Client{}}
}

type slackPayload struct {
	Text string `json:"text"`
}

// Notify posts a drift summary to a Slack incoming webhook.
func (s *slackNotifier) Notify(report *drift.Report) error {
	var msg string
	if !report.HasDrift() {
		msg = ":white_check_mark: *DriftWatch*: No infrastructure drift detected."
	} else {
		msg = fmt.Sprintf(":warning: *DriftWatch*: Drift detected in %d resource(s). Run `driftwatch` for details.", len(report.Drifts))
	}

	payload, err := json.Marshal(slackPayload{Text: msg})
	if err != nil {
		return fmt.Errorf("slack notify: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack notify: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
