package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/bus"
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
func (s *tenantStoreStub) GetTenantBySlug(context.Context, string) (*store.TenantData, error) { return nil, fmt.Errorf("not found") }
func (s *tenantStoreStub) ListTenants(context.Context) ([]store.TenantData, error) { return nil, nil }
func (s *tenantStoreStub) UpdateTenant(context.Context, uuid.UUID, map[string]any) error { return nil }
func (s *tenantStoreStub) AddUser(context.Context, uuid.UUID, string, string) error { return nil }
func (s *tenantStoreStub) RemoveUser(context.Context, uuid.UUID, string) error { return nil }
func (s *tenantStoreStub) GetUserRole(context.Context, uuid.UUID, string) (string, error) { return "", nil }
func (s *tenantStoreStub) ListUsers(context.Context, uuid.UUID) ([]store.TenantUserData, error) { return nil, nil }
func (s *tenantStoreStub) ListUserTenants(context.Context, string) ([]store.TenantUserData, error) { return nil, nil }
func (s *tenantStoreStub) GetTenantsByIDs(context.Context, []uuid.UUID) ([]store.TenantData, error) { return nil, nil }
func (s *tenantStoreStub) ResolveUserTenant(context.Context, string) (uuid.UUID, error) { return store.MasterTenantID, nil }
func (s *tenantStoreStub) GetTenantUser(context.Context, uuid.UUID) (*store.TenantUserData, error) { return nil, fmt.Errorf("not found") }
func (s *tenantStoreStub) CreateTenantUserReturning(context.Context, uuid.UUID, string, string, string) (*store.TenantUserData, error) {
	return &store.TenantUserData{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

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
