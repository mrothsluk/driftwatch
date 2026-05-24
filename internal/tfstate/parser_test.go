package tfstate_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/tfstate"
)

var sampleState = []byte(`{
  "version": 4,
  "resources": [
    {
      "type": "aws_instance",
      "name": "web",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "id": "i-0abc123",
            "instance_type": "t3.micro",
            "ami": "ami-0deadbeef"
          }
        }
      ]
    }
  ]
}`)

func TestParse_ValidState(t *testing.T) {
	state, err := tfstate.Parse(sampleState)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	res := state.Resources[0]
	if res.Type != "aws_instance" {
		t.Errorf("expected type aws_instance, got %q", res.Type)
	}
	if res.Name != "web" {
		t.Errorf("expected name web, got %q", res.Name)
	}
	if res.Attributes["id"] != "i-0abc123" {
		t.Errorf("expected id i-0abc123, got %v", res.Attributes["id"])
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := tfstate.Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParse_EmptyResources(t *testing.T) {
	data := []byte(`{"version": 4, "resources": []}`)
	state, err := tfstate.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(state.Resources))
	}
}
