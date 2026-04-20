package http

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestParseTenantSettingsDefaults(t *testing.T) {
	cfg := parseTenantSettings(nil)
	if cfg.WorkspaceMode != "isolated" {
		t.Fatalf("WorkspaceMode = %q, want isolated", cfg.WorkspaceMode)
	}
}

func TestTenantParentIDFromSettings(t *testing.T) {
	parent := uuid.New()
	raw, _ := json.Marshal(TenantSettings{ParentTenantID: parent.String()})
	if got := tenantParentID(raw); got != parent {
		t.Fatalf("tenantParentID() = %v, want %v", got, parent)
	}
}

func TestNormalizeReasoningOutputMode(t *testing.T) {
	cases := map[string]string{
		"FULL":    "full",
		"summary": "summary",
		"none":    "none",
		"bogus":   "full",
		"":        "full",
	}
	for in, want := range cases {
		if got := normalizeReasoningOutputMode(in); got != want {
			t.Fatalf("normalizeReasoningOutputMode(%q) = %q, want %q", in, got, want)
		}
	}
}
