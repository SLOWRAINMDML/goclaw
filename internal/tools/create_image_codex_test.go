package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nextlevelbuilder/goclaw/internal/providers"
)

type codexImageTokenSource struct {
	token string
}

func (s *codexImageTokenSource) Token() (string, error) { return s.token, nil }

func TestCreateImageTool_UsesCodexOAuthProvider(t *testing.T) {
	workspace := t.TempDir()
	wantPNG := []byte{0x89, 0x50, 0x4e, 0x47}
	b64 := base64.StdEncoding.EncodeToString(wantPNG)

	var capturedAuth string
	var capturedPath string
	var capturedBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		capturedPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &capturedBody)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{"b64_json": b64}},
		})
	}))
	defer srv.Close()

	reg := providers.NewRegistry(nil)
	reg.Register(providers.NewCodexProvider("codex-work", &codexImageTokenSource{token: "oauth-token"}, srv.URL, "gpt-image-1"))
	tool := NewCreateImageTool(reg)

	settings := BuiltinToolSettings{
		"create_image": []byte(`{"providers":[{"provider":"codex-work","model":"gpt-image-1","enabled":true,"timeout":30,"max_retries":1}]}`),
	}
	ctx := WithBuiltinToolSettings(context.Background(), settings)
	ctx = WithToolWorkspace(ctx, workspace)

	res := tool.Execute(ctx, map[string]any{"prompt": "draw a sunrise over mountains"})
	if res.IsError {
		t.Fatalf("Execute() error = %s", res.ForLLM)
	}
	if capturedAuth != "Bearer oauth-token" {
		t.Fatalf("Authorization = %q, want Bearer oauth-token", capturedAuth)
	}
	if capturedPath != "/images/generations" {
		t.Fatalf("path = %q, want /images/generations", capturedPath)
	}
	if got := capturedBody["model"]; got != "gpt-image-1" {
		t.Fatalf("model = %v, want gpt-image-1", got)
	}
	if len(res.Media) != 1 {
		t.Fatalf("media count = %d, want 1", len(res.Media))
	}
	if !strings.Contains(res.ForLLM, "MEDIA:") {
		t.Fatalf("ForLLM = %q, want MEDIA path", res.ForLLM)
	}
	if res.Provider != "codex-work" {
		t.Fatalf("Provider = %q, want codex-work", res.Provider)
	}
	if res.Model != "gpt-image-1" {
		t.Fatalf("Model = %q, want gpt-image-1", res.Model)
	}
}
