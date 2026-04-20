package http

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nextlevelbuilder/goclaw/internal/bus"
	"github.com/nextlevelbuilder/goclaw/internal/i18n"
	"github.com/nextlevelbuilder/goclaw/internal/skills"
	"github.com/nextlevelbuilder/goclaw/internal/store"
)

type skillFeedbackEntry struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	BehaviorRef      string `json:"behavior_ref,omitempty"`
	CodePath         string `json:"code_path,omitempty"`
	Prompt           string `json:"prompt,omitempty"`
	ExpectedBehavior string `json:"expected_behavior,omitempty"`
	ActualBehavior   string `json:"actual_behavior,omitempty"`
	Note             string `json:"note,omitempty"`
	CreatedBy        string `json:"created_by,omitempty"`
	CreatedAt        string `json:"created_at"`
}

func (h *SkillsHandler) feedbackFile(skillDir string) string {
	return filepath.Join(skillDir, "feedback", "examples.jsonl")
}

func (h *SkillsHandler) readSkillSource(skillDir string) (string, map[string]string, string, string, error) {
	contentBytes, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil {
		return "", nil, "", "", err
	}
	content := string(contentBytes)
	name, description, _, frontmatter := skills.ParseSkillFrontmatter(content)
	return content, frontmatter, name, description, nil
}

func copySkillTree(srcDir, destDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(destDir, 0755)
		}
		base := filepath.Base(rel)
		lowerRel := strings.ToLower(filepath.ToSlash(rel))
		if base == ".runtime" || strings.HasPrefix(lowerRel, ".runtime/") || strings.HasPrefix(lowerRel, "feedback/") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(lowerRel, "secrets/") || strings.HasPrefix(base, ".env") || strings.Contains(base, "credential") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(destDir, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		return err
	})
}

func (h *SkillsHandler) handleFork(w http.ResponseWriter, r *http.Request) {
	locale := store.LocaleFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": i18n.T(locale, i18n.MsgInvalidID, "skill")})
		return
	}
	info, ok := h.skills.GetSkillByID(r.Context(), id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": i18n.T(locale, i18n.MsgNotFound, "skill", id.String())})
		return
	}
	var req struct {
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		Visibility string `json:"visibility"`
	}
	if !bindJSON(w, r, locale, &req) {
		return
	}
	if req.Slug == "" {
		req.Slug = info.Slug + "-fork"
	}
	if !isValidSlug(req.Slug) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": i18n.T(locale, i18n.MsgInvalidSlug, "slug")})
		return
	}
	if req.Visibility == "" {
		req.Visibility = "internal"
	}
	skillDir := info.BaseDir
	if skillDir == "" {
		skillDir = filepath.Dir(info.Path)
	}
	content, frontmatter, parsedName, description, err := h.readSkillSource(skillDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if req.Name == "" {
		req.Name = parsedName + " Fork"
	}
	version := h.skills.GetNextVersion(r.Context(), req.Slug)
	if version <= 0 {
		version = 1
	}
	destDir := filepath.Join(h.tenantSkillsDir(r), req.Slug, fmt.Sprintf("%d", version))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := copySkillTree(skillDir, destDir); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	frontmatter["forked_from_skill_id"] = id.String()
	frontmatter["forked_from_slug"] = info.Slug
	frontmatter["forked_at"] = time.Now().UTC().Format(time.RFC3339)
	if runtimePrep, err := skills.PrepareRuntime(destDir, frontmatter); err == nil && runtimePrep != nil {
		_ = runtimePrep
	}
	if strings.TrimSpace(content) == "" {
		content = "---\nname: " + req.Name + "\n---\n"
	}
	desc := description
	createdID, err := h.skills.CreateSkillManaged(r.Context(), store.SkillCreateParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: &desc,
		OwnerID:     store.UserIDFromContext(r.Context()),
		Visibility:  req.Visibility,
		Status:      "active",
		Version:     version,
		FilePath:    destDir,
		FileSize:    int64(len(content)),
		Frontmatter: frontmatter,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.skills.BumpVersion()
	h.emitCacheInvalidate(bus.CacheKindSkills, createdID.String(), uuid.Nil)
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":      createdID,
		"slug":    req.Slug,
		"version": version,
		"path":    destDir,
	})
}

func (h *SkillsHandler) handlePrepareRuntime(w http.ResponseWriter, r *http.Request) {
	locale := store.LocaleFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": i18n.T(locale, i18n.MsgInvalidID, "skill")})
		return
	}
	info, ok := h.skills.GetSkillByID(r.Context(), id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": i18n.T(locale, i18n.MsgNotFound, "skill", id.String())})
		return
	}
	_, frontmatter, _, _, err := h.readSkillSource(info.BaseDir)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	result, err := skills.PrepareRuntime(info.BaseDir, frontmatter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"runtime": result})
}

func (h *SkillsHandler) handleListFeedback(w http.ResponseWriter, r *http.Request) {
	locale := store.LocaleFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": i18n.T(locale, i18n.MsgInvalidID, "skill")})
		return
	}
	info, ok := h.skills.GetSkillByID(r.Context(), id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": i18n.T(locale, i18n.MsgNotFound, "skill", id.String())})
		return
	}
	path := h.feedbackFile(info.BaseDir)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		writeJSON(w, http.StatusOK, map[string]any{"feedback": []skillFeedbackEntry{}})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer f.Close()
	entries := []skillFeedbackEntry{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var entry skillFeedbackEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err == nil {
			entries = append(entries, entry)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"feedback": entries})
}

func (h *SkillsHandler) handleAddFeedback(w http.ResponseWriter, r *http.Request) {
	locale := store.LocaleFromContext(r.Context())
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": i18n.T(locale, i18n.MsgInvalidID, "skill")})
		return
	}
	info, ok := h.skills.GetSkillByID(r.Context(), id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": i18n.T(locale, i18n.MsgNotFound, "skill", id.String())})
		return
	}
	var entry skillFeedbackEntry
	if !bindJSON(w, r, locale, &entry) {
		return
	}
	if entry.Type == "" {
		entry.Type = "issue"
	}
	entry.ID = uuid.NewString()
	entry.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	entry.CreatedBy = store.UserIDFromContext(r.Context())
	feedbackPath := h.feedbackFile(info.BaseDir)
	if err := os.MkdirAll(filepath.Dir(feedbackPath), 0755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	fh, err := os.OpenFile(feedbackPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer fh.Close()
	line, _ := json.Marshal(entry)
	if _, err := fh.Write(append(line, '\n')); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"feedback": entry})
}
