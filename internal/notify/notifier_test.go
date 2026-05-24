package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/notify"
)

func sampleReport(hasDrift bool) *drift.Report {
	r := &drift.Report{}
	if hasDrift {
		r.Drifts = []drift.ResourceDrift{
			{
				ResourceID:   "aws_instance.web",
				ResourceType: "aws_instance",
				Diffs: []drift.AttributeDiff{
					{Attribute: "instance_type", Expected: "t3.micro", Actual: "t3.small"},
				},
			},
		}
	}
	return r
}

func TestNew_UnknownKind(t *testing.T) {
	_, err := notify.New(notify.Config{Kind: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestNew_SlackMissingURL(t *testing.T) {
	_, err := notify.New(notify.Config{Kind: "slack", Options: map[string]string{}})
	if err == nil {
		t.Fatal("expected error when webhook_url missing")
	}
}

func TestNew_LogNotifier_NoDrift(t *testing.T) {
	n, err := notify.New(notify.Config{Kind: "log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(sampleReport(false)); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestNew_LogNotifier_WithDrift(t *testing.T) {
	n, err := notify.New(notify.Config{Kind: "log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(sampleReport(true)); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestWebhookNotifier_PostsJSON(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := notify.New(notify.Config{Kind: "webhook", Options: map[string]string{"url": ts.URL}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(sampleReport(true)); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if received["has_drift"] != true {
		t.Errorf("expected has_drift=true, got %v", received["has_drift"])
	}
}

func TestWebhookNotifier_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, err := notify.New(notify.Config{Kind: "webhook", Options: map[string]string{"url": ts.URL}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(sampleReport(false)); err == nil {
		t.Fatal("expected error on 500 response")
	}
}
