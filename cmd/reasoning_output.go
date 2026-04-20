package cmd

import (
	"context"
	"strings"

	"github.com/nextlevelbuilder/goclaw/internal/store"
)

func resolveTenantReasoningOutputMode(ctx context.Context, configs store.SystemConfigStore) string {
	if configs == nil {
		return "full"
	}
	value, err := configs.Get(ctx, "coding.reasoning_output")
	if err != nil {
		return "full"
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "summary", "none", "full":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "full"
	}
}
