package http

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/permissions"
	"github.com/nextlevelbuilder/goclaw/internal/store"
)

// TenantSettings stores optional hierarchical/workspace configuration inside tenants.settings.
// We keep this in JSON settings so the feature can ship without a schema migration.
type TenantSettings struct {
	ParentTenantID       string              `json:"parent_tenant_id,omitempty"`
	InheritParentAccess  bool                `json:"inherit_parent_access,omitempty"`
	WorkspaceMode        string              `json:"workspace_mode,omitempty"`
	CommonSkillIDs       []string            `json:"common_skill_ids,omitempty"`
	GitLinks             []TenantGitLink     `json:"git_links,omitempty"`
	ChannelBindings      []TenantChannelBind `json:"channel_bindings,omitempty"`
	OutputBindings       []TenantOutputBind  `json:"output_bindings,omitempty"`
	CodingProvider       string              `json:"coding_provider,omitempty"`
	CodingModel          string              `json:"coding_model,omitempty"`
	ReasoningOutput      string              `json:"reasoning_output,omitempty"`
}

type TenantGitLink struct {
	Name      string `json:"name,omitempty"`
	RepoURL   string `json:"repo_url,omitempty"`
	Branch    string `json:"branch,omitempty"`
	Directory string `json:"directory,omitempty"`
	Mode      string `json:"mode,omitempty"`
}

type TenantChannelBind struct {
	ChannelInstanceID string   `json:"channel_instance_id,omitempty"`
	MatchChatIDs      []string `json:"match_chat_ids,omitempty"`
	MatchTopics       []string `json:"match_topics,omitempty"`
	AgentID           string   `json:"agent_id,omitempty"`
	Enabled           bool     `json:"enabled,omitempty"`
}

type TenantOutputBind struct {
	Type       string            `json:"type,omitempty"`
	Target     string            `json:"target,omitempty"`
	Channel    string            `json:"channel,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func parseTenantSettings(raw json.RawMessage) TenantSettings {
	var cfg TenantSettings
	if len(raw) == 0 {
		return cfg
	}
	_ = json.Unmarshal(raw, &cfg)
	if cfg.WorkspaceMode == "" {
		cfg.WorkspaceMode = "isolated"
	}
	return cfg
}

func marshalTenantSettings(cfg TenantSettings) json.RawMessage {
	if cfg.WorkspaceMode == "" {
		cfg.WorkspaceMode = "isolated"
	}
	b, _ := json.Marshal(cfg)
	return b
}

func tenantParentID(raw json.RawMessage) uuid.UUID {
	cfg := parseTenantSettings(raw)
	if cfg.ParentTenantID == "" {
		return uuid.Nil
	}
	tid, err := uuid.Parse(cfg.ParentTenantID)
	if err != nil {
		return uuid.Nil
	}
	return tid
}

func listChildTenants(ctx context.Context, parentID uuid.UUID) ([]store.TenantData, error) {
	if pkgTenantCache == nil || pkgTenantCache.store == nil {
		return nil, fmt.Errorf("tenant store unavailable")
	}
	tenants, err := pkgTenantCache.store.ListTenants(ctx)
	if err != nil {
		return nil, err
	}
	children := make([]store.TenantData, 0)
	for _, tenant := range tenants {
		if tenantParentID(tenant.Settings) == parentID {
			children = append(children, tenant)
		}
	}
	return children, nil
}

func tenantRoleRank(role string) int {
	switch strings.ToLower(role) {
	case store.TenantRoleOwner:
		return 5
	case store.TenantRoleAdmin:
		return 4
	case store.TenantRoleOperator:
		return 3
	case store.TenantRoleMember:
		return 2
	case store.TenantRoleViewer:
		return 1
	default:
		return 0
	}
}

func tenantRoleMeetsMinimum(role string, min permissions.Role) bool {
	required := 1
	switch min {
	case permissions.RoleOwner:
		required = 5
	case permissions.RoleAdmin:
		required = 4
	case permissions.RoleOperator:
		required = 3
	case permissions.RoleViewer, "":
		required = 1
	}
	return tenantRoleRank(role) >= required
}

func userHasTenantOrAncestorAccess(ctx context.Context, tenantID uuid.UUID, userID string, min permissions.Role) bool {
	if tenantID == uuid.Nil || userID == "" || pkgTenantCache == nil || pkgTenantCache.store == nil {
		return false
	}
	visited := map[uuid.UUID]struct{}{}
	current := tenantID
	for current != uuid.Nil {
		if _, seen := visited[current]; seen {
			return false
		}
		visited[current] = struct{}{}
		role, err := pkgTenantCache.store.GetUserRole(ctx, current, userID)
		if err == nil && role != "" {
			if current == tenantID {
				return tenantRoleMeetsMinimum(role, min)
			}
			tenant, terr := pkgTenantCache.GetTenant(ctx, tenantID)
			if terr == nil && tenant != nil {
				cfg := parseTenantSettings(tenant.Settings)
				if cfg.InheritParentAccess && tenantRoleMeetsMinimum(role, permissions.RoleAdmin) {
					return true
				}
			}
		}
		tenant, err := pkgTenantCache.GetTenant(ctx, current)
		if err != nil || tenant == nil {
			return false
		}
		parentID := tenantParentID(tenant.Settings)
		current = parentID
	}
	return false
}

func normalizeReasoningOutputMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if slices.Contains([]string{"full", "summary", "none"}, mode) {
		return mode
	}
	return "full"
}
