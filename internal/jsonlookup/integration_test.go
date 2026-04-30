package jsonlookup_test

import (
	"encoding/json"
	"testing"

	"github.com/your-org/logpilot/internal/jsonlookup"
)

func TestLookup_ChainedApplication(t *testing.T) {
	regionTable := map[string]string{
		"us-east-1": "US East",
		"eu-west-1": "EU West",
	}
	envTable := map[string]string{
		"prod": "Production",
		"stg":  "Staging",
	}

	regionLookup, err := jsonlookup.New("region", regionTable, "region_label")
	if err != nil {
		t.Fatalf("failed to create region lookup: %v", err)
	}
	envLookup, err := jsonlookup.New("env", envTable, "env_label")
	if err != nil {
		t.Fatalf("failed to create env lookup: %v", err)
	}

	input := `{"region":"us-east-1","env":"prod","msg":"ok"}`
	out := regionLookup.Apply(input)
	out = envLookup.Apply(out)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["region_label"] != "US East" {
		t.Errorf("region_label: expected 'US East', got %v", obj["region_label"])
	}
	if obj["env_label"] != "Production" {
		t.Errorf("env_label: expected 'Production', got %v", obj["env_label"])
	}
}
