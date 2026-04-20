package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareRuntimePython(t *testing.T) {
	dir := t.TempDir()
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(scriptsDir, "run.py")
	if err := os.WriteFile(entry, []byte("print('ok')\n"), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := PrepareRuntime(dir, map[string]string{
		"runtime":    "python",
		"entrypoint": "scripts/run.py",
	})
	if err != nil {
		t.Fatalf("PrepareRuntime() error = %v", err)
	}
	if result == nil || !result.Prepared {
		t.Fatalf("PrepareRuntime() = %#v, want prepared result", result)
	}
	if _, err := os.Stat(filepath.Join(dir, ".runtime", "run.py")); err != nil {
		t.Fatalf("wrapper missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".runtime", "manifest.json")); err != nil {
		t.Fatalf("manifest missing: %v", err)
	}
}
