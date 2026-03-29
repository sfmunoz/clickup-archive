package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type GenerateContent struct {
	clickupDir string
	contentDir string
}

func NewGenerateContent(clickupDir, contentDir string) (*GenerateContent, error) {
	return &GenerateContent{clickupDir: clickupDir, contentDir: contentDir}, nil
}

func (g *GenerateContent) Run() error {
	entries, err := os.ReadDir(g.clickupDir)
	if err != nil {
		return fmt.Errorf("read clickup dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if err := g.processWorkspace(filepath.Join(g.clickupDir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

// pageDir returns the content output directory corresponding to a clickup entity directory.
func (g *GenerateContent) pageDir(entityDir string) string {
	rel, _ := filepath.Rel(g.clickupDir, entityDir)
	return filepath.Join(g.contentDir, rel)
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// childEntry holds the directory name (ID) and display name of a child entity.
type childEntry struct {
	id   string
	name string
}

// readChildren lists immediate subdirectories and reads each child's name from index.json.
func readChildren(dir string) ([]childEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var children []childEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		var named struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := readJSON(filepath.Join(dir, e.Name(), "index.json"), &named); err != nil {
			continue
		}
		children = append(children, childEntry{id: e.Name(), name: named.Name})
	}
	return children, nil
}

func writePage(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// yamlStr returns a safely double-quoted YAML string value.
func yamlStr(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

func (g *GenerateContent) processWorkspace(dir string) error {
	var w api.Workspace
	if err := readJSON(filepath.Join(dir, "index.json"), &w); err != nil {
		return fmt.Errorf("read workspace: %w", err)
	}
	log.Info("Workspace", "id", w.ID, "name", w.Name)

	children, err := readChildren(dir)
	if err != nil {
		return fmt.Errorf("list spaces: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "title: %s\n", yamlStr(w.Name))
	fmt.Fprintf(&sb, "clickup_id: %s\n", yamlStr(w.ID))
	fmt.Fprintf(&sb, "clickup_type: \"workspace\"\n")
	if w.Color != "" {
		fmt.Fprintf(&sb, "color: %s\n", yamlStr(w.Color))
	}
	if w.Avatar != "" {
		fmt.Fprintf(&sb, "avatar: %s\n", yamlStr(w.Avatar))
	}
	fmt.Fprintf(&sb, "---\n\n")

	if len(children) > 0 {
		fmt.Fprintf(&sb, "## Spaces\n\n")
		for _, c := range children {
			fmt.Fprintf(&sb, "- [%s](%s/)\n", c.name, c.id)
		}
		fmt.Fprintf(&sb, "\n")
	}

	if err := writePage(filepath.Join(g.pageDir(dir), "_index.md"), sb.String()); err != nil {
		return fmt.Errorf("write workspace page: %w", err)
	}

	for _, c := range children {
		if err := g.processSpace(filepath.Join(dir, c.id)); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenerateContent) processSpace(dir string) error {
	var s api.Space
	if err := readJSON(filepath.Join(dir, "index.json"), &s); err != nil {
		return fmt.Errorf("read space: %w", err)
	}
	log.Info("Space", "id", s.ID, "name", s.Name)

	children, err := readChildren(dir)
	if err != nil {
		return fmt.Errorf("list folders: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "title: %s\n", yamlStr(s.Name))
	fmt.Fprintf(&sb, "clickup_id: %s\n", yamlStr(s.ID))
	fmt.Fprintf(&sb, "clickup_type: \"space\"\n")
	fmt.Fprintf(&sb, "private: %t\n", s.Private)
	fmt.Fprintf(&sb, "archived: %t\n", s.Archived)
	fmt.Fprintf(&sb, "multiple_assignees: %t\n", s.MultipleAssignees)
	if len(s.Statuses) > 0 {
		fmt.Fprintf(&sb, "statuses:\n")
		for _, st := range s.Statuses {
			fmt.Fprintf(&sb, "  - status: %s\n", yamlStr(st.Status))
			fmt.Fprintf(&sb, "    type: %s\n", yamlStr(st.Type))
			if st.Color != "" {
				fmt.Fprintf(&sb, "    color: %s\n", yamlStr(st.Color))
			}
		}
	}
	fmt.Fprintf(&sb, "---\n\n")

	if len(children) > 0 {
		fmt.Fprintf(&sb, "## Folders\n\n")
		for _, c := range children {
			fmt.Fprintf(&sb, "- [%s](%s/)\n", c.name, c.id)
		}
		fmt.Fprintf(&sb, "\n")
	}

	if err := writePage(filepath.Join(g.pageDir(dir), "_index.md"), sb.String()); err != nil {
		return fmt.Errorf("write space page: %w", err)
	}

	for _, c := range children {
		if err := g.processFolder(filepath.Join(dir, c.id)); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenerateContent) processFolder(dir string) error {
	var f api.Folder
	if err := readJSON(filepath.Join(dir, "index.json"), &f); err != nil {
		return fmt.Errorf("read folder: %w", err)
	}
	log.Info("Folder", "id", f.ID, "name", f.Name)

	children, err := readChildren(dir)
	if err != nil {
		return fmt.Errorf("list lists: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "title: %s\n", yamlStr(f.Name))
	fmt.Fprintf(&sb, "clickup_id: %s\n", yamlStr(f.ID))
	fmt.Fprintf(&sb, "clickup_type: \"folder\"\n")
	fmt.Fprintf(&sb, "orderindex: %d\n", f.Orderindex)
	fmt.Fprintf(&sb, "hidden: %t\n", f.Hidden)
	fmt.Fprintf(&sb, "task_count: %s\n", yamlStr(f.TaskCount))
	fmt.Fprintf(&sb, "space_id: %s\n", yamlStr(f.Space.ID))
	fmt.Fprintf(&sb, "space_name: %s\n", yamlStr(f.Space.Name))
	fmt.Fprintf(&sb, "---\n\n")

	if len(children) > 0 {
		fmt.Fprintf(&sb, "## Lists\n\n")
		for _, c := range children {
			fmt.Fprintf(&sb, "- [%s](%s/)\n", c.name, c.id)
		}
		fmt.Fprintf(&sb, "\n")
	}

	if err := writePage(filepath.Join(g.pageDir(dir), "_index.md"), sb.String()); err != nil {
		return fmt.Errorf("write folder page: %w", err)
	}

	for _, c := range children {
		if err := g.processList(filepath.Join(dir, c.id)); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenerateContent) processList(dir string) error {
	var l api.List
	if err := readJSON(filepath.Join(dir, "index.json"), &l); err != nil {
		return fmt.Errorf("read list: %w", err)
	}
	log.Info("List", "id", l.ID, "name", l.Name)

	children, err := readChildren(dir)
	if err != nil {
		return fmt.Errorf("list tasks: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "title: %s\n", yamlStr(l.Name))
	fmt.Fprintf(&sb, "clickup_id: %s\n", yamlStr(l.ID))
	fmt.Fprintf(&sb, "clickup_type: \"list\"\n")
	fmt.Fprintf(&sb, "orderindex: %d\n", l.Orderindex)
	fmt.Fprintf(&sb, "task_count: %d\n", l.TaskCount)
	fmt.Fprintf(&sb, "archived: %t\n", l.Archived)
	fmt.Fprintf(&sb, "folder_id: %s\n", yamlStr(l.Folder.ID))
	fmt.Fprintf(&sb, "folder_name: %s\n", yamlStr(l.Folder.Name))
	fmt.Fprintf(&sb, "space_id: %s\n", yamlStr(l.Space.ID))
	fmt.Fprintf(&sb, "space_name: %s\n", yamlStr(l.Space.Name))
	fmt.Fprintf(&sb, "---\n\n")

	if len(children) > 0 {
		fmt.Fprintf(&sb, "## Tasks\n\n")
		for _, c := range children {
			fmt.Fprintf(&sb, "- [%s](%s/)\n", c.name, c.id)
		}
		fmt.Fprintf(&sb, "\n")
	}

	if err := writePage(filepath.Join(g.pageDir(dir), "_index.md"), sb.String()); err != nil {
		return fmt.Errorf("write list page: %w", err)
	}

	for _, c := range children {
		if err := g.processTask(filepath.Join(dir, c.id)); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenerateContent) processTask(dir string) error {
	var t api.Task
	if err := readJSON(filepath.Join(dir, "index.json"), &t); err != nil {
		return fmt.Errorf("read task: %w", err)
	}
	log.Info("Task", "id", t.ID, "name", t.Name)

	var sb strings.Builder
	fmt.Fprintf(&sb, "---\n")
	fmt.Fprintf(&sb, "title: %s\n", yamlStr(t.Name))
	fmt.Fprintf(&sb, "clickup_id: %s\n", yamlStr(t.ID))
	fmt.Fprintf(&sb, "clickup_type: \"task\"\n")
	if len(t.Subtasks) > 0 {
		fmt.Fprintf(&sb, "subtasks:\n")
		for _, sub := range t.Subtasks {
			fmt.Fprintf(&sb, "  - id: %s\n", yamlStr(sub.ID))
			fmt.Fprintf(&sb, "    name: %s\n", yamlStr(sub.Name))
		}
	}
	fmt.Fprintf(&sb, "---\n")

	if err := writePage(filepath.Join(g.pageDir(dir), "index.md"), sb.String()); err != nil {
		return fmt.Errorf("write task page: %w", err)
	}
	return nil
}
