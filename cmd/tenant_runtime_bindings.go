package cmd

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/bus"
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

func resolveAgentRouteForMessage(ctx context.Context, deps *ConsumerDeps, msg bus.InboundMessage) string {
	if deps != nil && deps.TenantStore != nil && msg.TenantID != uuid.Nil {
		if tenant, err := deps.TenantStore.GetTenant(ctx, msg.TenantID); err == nil && tenant != nil {
			settings := parseTenantSettingsForRuntime(tenant.Settings)
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
		}
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
