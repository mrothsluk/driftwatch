// Package notify provides notification backends for reporting drift results.
package notify

import (
	"fmt"

	"github.com/example/driftwatch/internal/drift"
)

// Notifier sends drift reports to an external destination.
type Notifier interface {
	Notify(report *drift.Report) error
}

// Config holds configuration for a notifier.
type Config struct {
	Kind    string            // "slack", "webhook", "log"
	Options map[string]string // backend-specific options
}

// New constructs a Notifier from the given Config.
func New(cfg Config) (Notifier, error) {
	switch cfg.Kind {
	case "slack":
		webhookURL, ok := cfg.Options["webhook_url"]
		if !ok || webhookURL == "" {
			return nil, fmt.Errorf("notify: slack requires 'webhook_url' option")
		}
		return newSlackNotifier(webhookURL), nil
	case "webhook":
		url, ok := cfg.Options["url"]
		if !ok || url == "" {
			return nil, fmt.Errorf("notify: webhook requires 'url' option")
		}
		return newWebhookNotifier(url), nil
	case "log":
		return &logNotifier{}, nil
	default:
		return nil, fmt.Errorf("notify: unsupported kind %q", cfg.Kind)
	}
}
