package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/bus"
	"github.com/nextlevelbuilder/goclaw/internal/providers"
	"github.com/nextlevelbuilder/goclaw/internal/store"
)

type tenantStoreStub struct {
	tenant *store.TenantData
}

func (s *tenantStoreStub) CreateTenant(context.Context, *store.TenantData) error { return nil }
func (s *tenantStoreStub) GetTenant(ctx context.Context, id uuid.UUID) (*store.TenantData, error) {
	if s.tenant != nil && s.tenant.ID == id {
		return s.tenant, nil
	}
	return nil, fmt.Errorf("not found")
}
func (s *tenantStoreStub) GetTenantBySlug(context.Context, string) (*store.TenantData, error) {
	return nil, fmt.Errorf("not found")
}
func (s *tenantStoreStub) ListTenants(context.Context) ([]store.TenantData, error) { return nil, nil }
func (s *tenantStoreStub) UpdateTenant(context.Context, uuid.UUID, map[string]any) error {
	return nil
}
func (s *tenantStoreStub) AddUser(context.Context, uuid.UUID, string, string) error { return nil }
func (s *tenantStoreStub) RemoveUser(context.Context, uuid.UUID, string) error { return nil }
func (s *tenantStoreStub) GetUserRole(context.Context, uuid.UUID, string) (string, error) {
	return "", nil
}
func (s *tenantStoreStub) ListUsers(context.Context, uuid.UUID) ([]store.TenantUserData, error) {
	return nil, nil
}
func (s *tenantStoreStub) ListUserTenants(context.Context, string) ([]store.TenantUserData, error) {
	return nil, nil
}
func (s *tenantStoreStub) GetTenantsByIDs(context.Context, []uuid.UUID) ([]store.TenantData, error) {
	return nil, nil
}
func (s *tenantStoreStub) ResolveUserTenant(context.Context, string) (uuid.UUID, error) {
	return store.MasterTenantID, nil
}
func (s *tenantStoreStub) GetTenantUser(context.Context, uuid.UUID) (*store.TenantUserData, error) {
	return nil, fmt.Errorf("not found")
}
func (s *tenantStoreStub) CreateTenantUserReturning(context.Context, uuid.UUID, string, string, string) (*store.TenantUserData, error) {
	return &store.TenantUserData{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

type tenantRuntimeTestProvider struct {
	name         string
	defaultModel string
}

func (p *tenantRuntimeTestProvider) Chat(context.Context, providers.ChatRequest) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{}, nil
}

func (p *tenantRuntimeTestProvider) ChatStream(context.Context, providers.ChatRequest, func(providers.StreamChunk)) (*providers.ChatResponse, error) {
	return &providers.ChatResponse{}, nil
}

func (p *tenantRuntimeTestProvider) DefaultModel() string { return p.defaultModel }
func (p *tenantRuntimeTestProvider) Name() string         { return p.name }

func TestResolveAgentRouteForMessageUsesTenantChannelBinding(t *testing.T) {
	tenantID := uuid.New()
	settings, _ := json.Marshal(tenantSettings{
		ChannelBindings: []tenantChannelBinding{{
			ChannelInstanceID: "discord-main",
			MatchChatIDs:      []string{"chat-1"},
			MatchTopics:       []string{"42"},
			AgentID:           "agent-special",
			Enabled:           true,
		}},
	})
	deps := &ConsumerDeps{
		Cfg:         nil,
		TenantStore: &tenantStoreStub{tenant: &store.TenantData{ID: tenantID, Settings: settings}},
	}
	msg := bus.InboundMessage{
		TenantID: tenantID,
		Channel:  "discord-main",
		ChatID:   "chat-1",
		PeerKind: "group",
		Metadata: map[string]string{"message_thread_id": "42"},
	}
	if got := resolveAgentRouteForMessage(context.Background(), deps, msg); got != "agent-special" {
		t.Fatalf("resolveAgentRouteForMessage() = %q, want agent-special", got)
	}
}

func TestResolveTenantCodingOverrideUsesTenantProviderAndModel(t *testing.T) {
	tenantID := uuid.New()
	providerReg := providers.NewRegistry(nil)
	providerReg.RegisterForTenant(tenantID, &tenantRuntimeTestProvider{name: "codex-work", defaultModel: "gpt-5-codex"})

	settings, _ := json.Marshal(tenantSettings{
		CodingProvider: "codex-work",
		CodingModel:    "gpt-5.4-codex",
	})
	deps := &ConsumerDeps{
		TenantStore:      &tenantStoreStub{tenant: &store.TenantData{ID: tenantID, Settings: settings}},
		ProviderRegistry: providerReg,
	}

	providerOverride, modelOverride := resolveTenantCodingOverride(context.Background(), deps, tenantID)
	if providerOverride == nil {
		t.Fatal("providerOverride = nil, want tenant provider")
	}
	if providerOverride.Name() != "codex-work" {
		t.Fatalf("providerOverride.Name() = %q, want codex-work", providerOverride.Name())
	}
	if modelOverride != "gpt-5.4-codex" {
		t.Fatalf("modelOverride = %q, want gpt-5.4-codex", modelOverride)
	}
}

func TestResolveTenantCodingOverrideFallsBackToProviderDefaultModel(t *testing.T) {
	tenantID := uuid.New()
	providerReg := providers.NewRegistry(nil)
	providerReg.RegisterForTenant(tenantID, &tenantRuntimeTestProvider{name: "codex-work", defaultModel: "gpt-5-codex"})

	settings, _ := json.Marshal(tenantSettings{CodingProvider: "codex-work"})
	deps := &ConsumerDeps{
		TenantStore:      &tenantStoreStub{tenant: &store.TenantData{ID: tenantID, Settings: settings}},
		ProviderRegistry: providerReg,
	}

	providerOverride, modelOverride := resolveTenantCodingOverride(context.Background(), deps, tenantID)
	if providerOverride == nil {
		t.Fatal("providerOverride = nil, want tenant provider")
	}
	if modelOverride != "gpt-5-codex" {
		t.Fatalf("modelOverride = %q, want gpt-5-codex", modelOverride)
	}
}

func TestResolveTenantCodingOverrideModelOnlyLeavesProviderUnset(t *testing.T) {
	tenantID := uuid.New()
	settings, _ := json.Marshal(tenantSettings{CodingModel: "gpt-5-mini"})
	deps := &ConsumerDeps{
		TenantStore: &tenantStoreStub{tenant: &store.TenantData{ID: tenantID, Settings: settings}},
	}

	providerOverride, modelOverride := resolveTenantCodingOverride(context.Background(), deps, tenantID)
	if providerOverride != nil {
		t.Fatalf("providerOverride = %v, want nil", providerOverride)
	}
	if modelOverride != "gpt-5-mini" {
		t.Fatalf("modelOverride = %q, want gpt-5-mini", modelOverride)
	}
}

func TestResolveTenantCodingOverrideIgnoresMissingProviderOverride(t *testing.T) {
	tenantID := uuid.New()
	settings, _ := json.Marshal(tenantSettings{
		CodingProvider: "missing-provider",
		CodingModel:    "gpt-5-codex",
	})
	deps := &ConsumerDeps{
		TenantStore:      &tenantStoreStub{tenant: &store.TenantData{ID: tenantID, Settings: settings}},
		ProviderRegistry: providers.NewRegistry(nil),
	}

	providerOverride, modelOverride := resolveTenantCodingOverride(context.Background(), deps, tenantID)
	if providerOverride != nil {
		t.Fatalf("providerOverride = %v, want nil", providerOverride)
	}
	if modelOverride != "" {
		t.Fatalf("modelOverride = %q, want empty when provider lookup fails", modelOverride)
	}
}
