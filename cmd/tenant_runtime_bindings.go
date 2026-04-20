package cmd

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/agent"
	"github.com/nextlevelbuilder/goclaw/internal/bus"
	"github.com/nextlevelbuilder/goclaw/internal/providers"
	"github.com/nextlevelbuilder/goclaw/internal/store"
	"github.com/nextlevelbuilder/goclaw/internal/tools"
)

type tenantSettings struct {
	ChannelBindings []tenantChannelBinding `json:"channel_bindings,omitempty"`
	OutputBindings  []tenantOutputBinding  `json:"output_bindings,omitempty"`
	CodingProvider  string                 `json:"coding_provider,omitempty"`
	CodingModel     string                 `json:"coding_model,omitempty"`
}

type tenantChannelBinding struct {
	ChannelInstanceID string   `json:"channel_instance_id,omitempty"`
	MatchChatIDs      []string `json:"match_chat_ids,omitempty"`
	MatchTopics       []string `json:"match_topics,omitempty"`
	AgentID           string   `json:"agent_id,omitempty"`
	Enabled           bool     `json:"enabled,omitempty"`
}

type tenantOutputBinding struct {
	Type       string            `json:"type,omitempty"`
	Target     string            `json:"target,omitempty"`
	Channel    string            `json:"channel,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func parseTenantSettingsForRuntime(raw json.RawMessage) tenantSettings {
	var cfg tenantSettings
	if len(raw) == 0 {
		return cfg
	}
	_ = json.Unmarshal(raw, &cfg)
	return cfg
}

func loadTenantRuntimeSettings(ctx context.Context, deps *ConsumerDeps, tenantID uuid.UUID) tenantSettings {
	if deps == nil {
		return tenantSettings{}
	}
	return loadTenantRuntimeSettingsFromStore(ctx, deps.TenantStore, tenantID)
}

func loadTenantRuntimeSettingsFromStore(ctx context.Context, tenantStore store.TenantStore, tenantID uuid.UUID) tenantSettings {
	if tenantStore == nil || tenantID == uuid.Nil {
		return tenantSettings{}
	}
	tenant, err := tenantStore.GetTenant(ctx, tenantID)
	if err != nil || tenant == nil {
		return tenantSettings{}
	}
	return parseTenantSettingsForRuntime(tenant.Settings)
}

func resolveTenantCodingOverride(ctx context.Context, deps *ConsumerDeps, tenantID uuid.UUID) (providers.Provider, string) {
	if deps == nil {
		return nil, ""
	}
	return resolveTenantCodingOverrideFromStore(ctx, deps.TenantStore, deps.ProviderRegistry, tenantID)
}

func resolveTenantCodingOverrideFromStore(ctx context.Context, tenantStore store.TenantStore, providerRegistry *providers.Registry, tenantID uuid.UUID) (providers.Provider, string) {
	settings := loadTenantRuntimeSettingsFromStore(ctx, tenantStore, tenantID)
	providerName := strings.TrimSpace(settings.CodingProvider)
	modelName := strings.TrimSpace(settings.CodingModel)
	if providerName == "" {
		return nil, modelName
	}
	if providerRegistry == nil {
		slog.Warn("tenant.coding_override.registry_missing", "tenant_id", tenantID, "provider", providerName)
		return nil, ""
	}
	providerOverride, err := providerRegistry.GetForTenant(tenantID, providerName)
	if err != nil || providerOverride == nil {
		slog.Warn("tenant.coding_override.provider_not_found", "tenant_id", tenantID, "provider", providerName, "error", err)
		return nil, ""
	}
	if modelName == "" {
		modelName = providerOverride.DefaultModel()
	}
	return providerOverride, modelName
}

func applyTenantCodingOverride(ctx context.Context, tenantStore store.TenantStore, providerRegistry *providers.Registry, tenantID uuid.UUID, req *agent.RunRequest) {
	if req == nil {
		return
	}
	providerOverride, modelOverride := resolveTenantCodingOverrideFromStore(ctx, tenantStore, providerRegistry, tenantID)
	if providerOverride != nil {
		req.ProviderOverride = providerOverride
	}
	if modelOverride != "" {
		req.ModelOverride = modelOverride
	}
}

func resolveAgentRouteForMessage(ctx context.Context, deps *ConsumerDeps, msg bus.InboundMessage) string {
	settings := loadTenantRuntimeSettings(ctx, deps, msg.TenantID)
	topicID := strings.TrimSpace(msg.Metadata[tools.MetaMessageThreadID])
	for _, binding := range settings.ChannelBindings {
		if !binding.Enabled || binding.ChannelInstanceID == "" || binding.AgentID == "" {
			continue
		}
		if binding.ChannelInstanceID != msg.Channel {
			continue
		}
		if len(binding.MatchChatIDs) > 0 && !containsString(binding.MatchChatIDs, msg.ChatID) {
			continue
		}
		if len(binding.MatchTopics) > 0 && !containsString(binding.MatchTopics, topicID) {
			continue
		}
		return binding.AgentID
	}
	return resolveAgentRoute(deps.Cfg, msg.Channel, msg.ChatID, msg.PeerKind)
}

func containsString(items []string, want string) bool {
	for _, item := range items {
		if strings.TrimSpace(item) == strings.TrimSpace(want) {
			return true
		}
	}
	return false
}
