package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RuntimePrepareResult describes generated runtime artifacts for a skill.
type RuntimePrepareResult struct {
	Runtime   string `json:"runtime,omitempty"`
	Entrypoint string `json:"entrypoint,omitempty"`
	Wrapper   string `json:"wrapper,omitempty"`
	Manifest  string `json:"manifest,omitempty"`
	Prepared  bool   `json:"prepared"`
}

// PrepareRuntime materializes lightweight execution artifacts under <skillDir>/.runtime.
// Current MVP focuses on Python skills declared via frontmatter or a scripts/*.py entrypoint.
func PrepareRuntime(skillDir string, frontmatter map[string]string) (*RuntimePrepareResult, error) {
	runtimeName := strings.TrimSpace(frontmatter["runtime"])
	entrypoint := strings.TrimSpace(frontmatter["entrypoint"])
	if runtimeName == "" {
		if entrypoint == "" {
			candidate := filepath.Join(skillDir, "scripts", "run.py")
			if _, err := os.Stat(candidate); err == nil {
				entrypoint = "scripts/run.py"
			}
		}
		if strings.HasSuffix(strings.ToLower(entrypoint), ".py") {
			runtimeName = "python"
		}
	}
	if runtimeName != "python" || entrypoint == "" {
		return &RuntimePrepareResult{Prepared: false, Runtime: runtimeName, Entrypoint: entrypoint}, nil
	}
	absEntry := filepath.Join(skillDir, filepath.Clean(entrypoint))
	if _, err := os.Stat(absEntry); err != nil {
		return nil, fmt.Errorf("runtime entrypoint missing: %s", entrypoint)
	}
	runtimeDir := filepath.Join(skillDir, ".runtime")
	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		return nil, err
	}
	wrapperPath := filepath.Join(runtimeDir, "run.py")
	manifestPath := filepath.Join(runtimeDir, "manifest.json")
	wrapper := fmt.Sprintf("import runpy\nimport pathlib\nBASE = pathlib.Path(__file__).resolve().parent.parent\nTARGET = BASE / %q\nrunpy.run_path(str(TARGET), run_name='__main__')\n", filepath.ToSlash(entrypoint))
	if err := os.WriteFile(wrapperPath, []byte(wrapper), 0644); err != nil {
		return nil, err
	}
	manifest := map[string]any{
		"prepared_at": time.Now().UTC().Format(time.RFC3339),
		"runtime":     "python",
		"entrypoint":  filepath.ToSlash(entrypoint),
		"wrapper":     filepath.ToSlash(filepath.Base(wrapperPath)),
		"skill_dir":   skillDir,
	}
	manifestBytes, _ := json.MarshalIndent(manifest, "", "  ")
	if err := os.WriteFile(manifestPath, manifestBytes, 0644); err != nil {
		return nil, err
	}
	// Best-effort compile. Ignore failures so runtime prep still succeeds on minimal environments.
	_, _ = exec.Command("python3", "-m", "compileall", runtimeDir).CombinedOutput()
	return &RuntimePrepareResult{
		Runtime:    "python",
		Entrypoint: filepath.ToSlash(entrypoint),
		Wrapper:    wrapperPath,
		Manifest:   manifestPath,
		Prepared:   true,
	}, nil
}
