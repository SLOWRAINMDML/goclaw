package channels

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/bus"
	"github.com/nextlevelbuilder/goclaw/internal/store"
)

type tenantRuntimeSettings struct {
	OutputBindings []tenantRuntimeOutputBinding `json:"output_bindings,omitempty"`
}

type tenantRuntimeOutputBinding struct {
	Type       string            `json:"type,omitempty"`
	Target     string            `json:"target,omitempty"`
	Channel    string            `json:"channel,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func parseTenantRuntimeSettings(raw json.RawMessage) tenantRuntimeSettings {
	var cfg tenantRuntimeSettings
	if len(raw) == 0 {
		return cfg
	}
	_ = json.Unmarshal(raw, &cfg)
	return cfg
}

func (m *Manager) expandOutboundBindings(ctx context.Context, msg bus.OutboundMessage) []bus.OutboundMessage {
	if m == nil || m.tenantStore == nil || msg.TenantID == uuid.Nil {
		return nil
	}
	if msg.Metadata != nil && msg.Metadata["_output_binding_applied"] == "true" {
		return nil
	}
	tenantCtx := store.WithCrossTenant(context.Background())
	tenantCtx = store.WithTenantID(tenantCtx, msg.TenantID)
	tenant, err := m.tenantStore.GetTenant(tenantCtx, msg.TenantID)
	if err != nil || tenant == nil {
		return nil
	}
	settings := parseTenantRuntimeSettings(tenant.Settings)
	outs := make([]bus.OutboundMessage, 0)
	for _, binding := range settings.OutputBindings {
		bindingType := strings.ToLower(strings.TrimSpace(binding.Type))
		if bindingType == "" {
			bindingType = "channel_copy"
		}
		if bindingType != "channel_copy" && bindingType != "api_channel" && bindingType != "channel" {
			continue
		}
		if strings.TrimSpace(binding.Channel) == "" || strings.TrimSpace(binding.Target) == "" {
			continue
		}
		clone := msg
		clone.Channel = strings.TrimSpace(binding.Channel)
		clone.ChatID = strings.TrimSpace(binding.Target)
		clone.Metadata = copyStringMap(msg.Metadata)
		if clone.Metadata == nil {
			clone.Metadata = map[string]string{}
		}
		clone.Metadata["_output_binding_applied"] = "true"
		for k, v := range binding.Parameters {
			if strings.TrimSpace(k) != "" && v != "" {
				clone.Metadata[k] = v
			}
		}
		outs = append(outs, clone)
	}
	return outs
}

func copyStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
