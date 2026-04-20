package channels

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
func (s *tenantStoreStub) GetTenant(_ context.Context, id uuid.UUID) (*store.TenantData, error) {
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

func TestExpandOutboundBindings(t *testing.T) {
	tenantID := uuid.New()
	settings, _ := json.Marshal(tenantRuntimeSettings{
		OutputBindings: []tenantRuntimeOutputBinding{{
			Type:    "channel_copy",
			Channel: "discord-secondary",
			Target:  "chat-2",
			Parameters: map[string]string{
				"message_thread_id": "77",
			},
		}},
	})
	mgr := NewManager(nil)
	mgr.SetTenantStore(&tenantStoreStub{tenant: &store.TenantData{ID: tenantID, Settings: settings}})
	outs := mgr.expandOutboundBindings(context.Background(), bus.OutboundMessage{
		Channel:  "discord-main",
		ChatID:   "chat-1",
		Content:  "hello",
		TenantID: tenantID,
	})
	if len(outs) != 1 {
		t.Fatalf("expandOutboundBindings() len = %d, want 1", len(outs))
	}
	if outs[0].Channel != "discord-secondary" || outs[0].ChatID != "chat-2" {
		t.Fatalf("unexpected clone target: %#v", outs[0])
	}
	if outs[0].Metadata["_output_binding_applied"] != "true" {
		t.Fatalf("expected clone metadata flag")
	}
	if outs[0].Metadata["message_thread_id"] != "77" {
		t.Fatalf("expected parameter copy, got %#v", outs[0].Metadata)
	}
}
