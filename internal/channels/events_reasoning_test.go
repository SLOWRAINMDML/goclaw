package channels

import "testing"

func TestFormatReasoningPreviewForMode(t *testing.T) {
	if got := formatReasoningPreviewForMode("none", "secret reasoning"); got != "" {
		t.Fatalf("none mode = %q, want empty", got)
	}
	summary := formatReasoningPreviewForMode("summary", "first second third")
	if summary == "" || summary == "_Reasoning:_\nfirst second third" {
		t.Fatalf("summary mode did not produce compact preview: %q", summary)
	}
	full := formatReasoningPreviewForMode("full", "first second third")
	if full == "" {
		t.Fatalf("full mode returned empty")
	}
}
