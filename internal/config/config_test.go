package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_ValidMinimal(t *testing.T) {
	yaml := []byte(`statefile: terraform.tfstate
provider: aws
`)
	cfg, err := Parse(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Statefile != "terraform.tfstate" {
		t.Errorf("statefile: got %q, want %q", cfg.Statefile, "terraform.tfstate")
	}
	if cfg.Provider != ProviderAWS {
		t.Errorf("provider: got %q, want %q", cfg.Provider, ProviderAWS)
	}
	if cfg.Output != "text" {
		t.Errorf("output default: got %q, want \"text\"", cfg.Output)
	}
}

func TestParse_AllFields(t *testing.T) {
	yaml := []byte(`
statefile: state.json
provider: gcp
output: json
filters:
  resource_types:
    - google_compute_instance
  exclude_ids:
    - projects/my-proj/zones/us-central1-a/instances/old-vm
`)
	cfg, err := Parse(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Output != "json" {
		t.Errorf("output: got %q", cfg.Output)
	}
	if len(cfg.Filters.ResourceTypes) != 1 {
		t.Errorf("resource_types length: got %d", len(cfg.Filters.ResourceTypes))
	}
	if len(cfg.Filters.ExcludeIDs) != 1 {
		t.Errorf("exclude_ids length: got %d", len(cfg.Filters.ExcludeIDs))
	}
}

func TestParse_MissingStatefile(t *testing.T) {
	_, err := Parse([]byte(`provider: aws
`))
	if err == nil {
		t.Fatal("expected error for missing statefile")
	}
}

func TestParse_MissingProvider(t *testing.T) {
	_, err := Parse([]byte(`statefile: s.tfstate
`))
	if err == nil {
		t.Fatal("expected error for missing provider")
	}
}

func TestParse_UnsupportedProvider(t *testing.T) {
	_, err := Parse([]byte(`statefile: s.tfstate
provider: digitalocean
`))
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestParse_InvalidOutput(t *testing.T) {
	_, err := Parse([]byte(`statefile: s.tfstate
provider: azure
output: xml
`))
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	content := []byte(`statefile: prod.tfstate
provider: gcp
`)
	if err := os.WriteFile(p, content, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Provider != ProviderGCP {
		t.Errorf("provider: got %q", cfg.Provider)
	}
}
