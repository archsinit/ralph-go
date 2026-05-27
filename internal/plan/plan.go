package plan

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TaskSpec defines a single task in a plan.
type TaskSpec struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

// Plan represents a checklist/plan with ordered tasks.
type Plan struct {
	Tasks []TaskSpec
}

// Decode parses a JSON plan response from an agent.
func Decode(jsonStr string) (*Plan, error) {
	// Try to parse as JSON array of tasks
	var tasks []TaskSpec
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&tasks); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate that each task has required fields
	for i, task := range tasks {
		if task.ID == "" {
			return nil, fmt.Errorf("task %d missing required field 'id'", i)
		}
		if task.Title == "" {
			return nil, fmt.Errorf("task %d missing required field 'title'", i)
		}
	}

	return &Plan{Tasks: tasks}, nil
}

// Write writes the plan to files in the output directory.
// Returns the list of files written.
func (p *Plan) Write(outDir string) ([]string, error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	var written []string

	// Write plan.md (index file)
	planMDPath := filepath.Join(outDir, "plan.md")
	planContent := p.renderPlanMD()
	if err := os.WriteFile(planMDPath, []byte(planContent), 0644); err != nil {
		return nil, fmt.Errorf("write plan.md: %w", err)
	}
	written = append(written, planMDPath)

	// Write task files in tasks/ directory
	tasksDir := filepath.Join(outDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return nil, fmt.Errorf("create tasks directory: %w", err)
	}

	for i, task := range p.Tasks {
		taskNum := i + 1
		filename := fmt.Sprintf("%02d-%s.md", taskNum, slugify(task.ID))
		filepath := filepath.Join(tasksDir, filename)

		content := fmt.Sprintf("# %s\n\n%s\n", task.Title, task.Description)
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("write task file %s: %w", filename, err)
		}
		written = append(written, filepath)
	}

	return written, nil
}

// renderPlanMD generates the plan.md index file content.
func (p *Plan) renderPlanMD() string {
	var sb strings.Builder
	sb.WriteString("# Plan\n\n")

	for i, task := range p.Tasks {
		taskNum := i + 1
		filename := fmt.Sprintf("%02d-%s.md", taskNum, slugify(task.ID))
		sb.WriteString(fmt.Sprintf("- [%s](%s) — %s\n", task.Title, filepath.Join("tasks", filename), task.Description))
	}

	return sb.String()
}

// slugify converts a string to a URL-safe slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result.WriteRune(ch)
		}
	}
	s = result.String()
	// Clean up multiple consecutive hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	return s
}
