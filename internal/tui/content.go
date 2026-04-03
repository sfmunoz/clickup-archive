package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/sfmunoz/clickup-archive/internal/archive"
)

func renderContent(node any, width, _ int) string {
	switch n := node.(type) {
	case *archive.Workspace:
		return renderWorkspace(n, width)
	case *archive.Space:
		return renderSpace(n, width)
	case *archive.Folder:
		return renderFolder(n)
	case *archive.List:
		return renderList(n)
	case *archive.Task:
		return renderTask(n, width)
	case *archive.Comment:
		return renderComment(n, width)
	default:
		return "Select an item to view details"
	}
}

func renderWorkspace(w *archive.Workspace, _ int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Workspace\n\n")
	fmt.Fprintf(&b, "Name:    %s\n", w.Data.Name)
	fmt.Fprintf(&b, "ID:      %s\n", w.Data.ID)
	fmt.Fprintf(&b, "Members: %d\n", len(w.Data.Members))
	return b.String()
}

func renderSpace(s *archive.Space, _ int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Space\n\n")
	fmt.Fprintf(&b, "Name:     %s\n", s.Data.Name)
	fmt.Fprintf(&b, "ID:       %s\n", s.Data.ID)
	fmt.Fprintf(&b, "Private:  %v\n", s.Data.Private)
	fmt.Fprintf(&b, "Archived: %v\n", s.Data.Archived)
	if len(s.Data.Statuses) > 0 {
		names := make([]string, len(s.Data.Statuses))
		for i, st := range s.Data.Statuses {
			names[i] = st.Status
		}
		fmt.Fprintf(&b, "Statuses: %s\n", strings.Join(names, ", "))
	}
	return b.String()
}

func renderFolder(f *archive.Folder) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Folder\n\n")
	fmt.Fprintf(&b, "Name:       %s\n", f.Data.Name)
	fmt.Fprintf(&b, "ID:         %s\n", f.Data.ID)
	fmt.Fprintf(&b, "Task count: %s\n", f.Data.TaskCount)
	fmt.Fprintf(&b, "Hidden:     %v\n", f.Data.Hidden)
	return b.String()
}

func renderList(l *archive.List) string {
	var b strings.Builder
	fmt.Fprintf(&b, "List\n\n")
	fmt.Fprintf(&b, "Name:       %s\n", l.Data.Name)
	fmt.Fprintf(&b, "ID:         %s\n", l.Data.ID)
	fmt.Fprintf(&b, "Task count: %d\n", l.Data.TaskCount)
	fmt.Fprintf(&b, "Archived:   %v\n", l.Data.Archived)
	return b.String()
}

func renderTask(t *archive.Task, width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Task\n\n")
	fmt.Fprintf(&b, "Name:      %s\n", t.Data.Name)
	fmt.Fprintf(&b, "ID:        %s\n", t.Data.ID)
	fmt.Fprintf(&b, "Status:    %s\n", t.Data.Status.Status)
	if t.Data.Priority != nil {
		fmt.Fprintf(&b, "Priority:  %s\n", t.Data.Priority.Priority)
	} else {
		fmt.Fprintf(&b, "Priority:  —\n")
	}
	if len(t.Data.Assignees) > 0 {
		names := make([]string, len(t.Data.Assignees))
		for i, a := range t.Data.Assignees {
			names[i] = a.Username
		}
		fmt.Fprintf(&b, "Assignees: %s\n", strings.Join(names, ", "))
	} else {
		fmt.Fprintf(&b, "Assignees: —\n")
	}
	if t.Data.DueDate != nil {
		fmt.Fprintf(&b, "Due date:  %s\n", formatEpochMs(*t.Data.DueDate))
	} else {
		fmt.Fprintf(&b, "Due date:  —\n")
	}
	if t.Data.TextContent != "" {
		fmt.Fprintf(&b, "\n")
		maxW := max(width-4, 20)
		fmt.Fprintf(&b, "%s\n", truncate(t.Data.TextContent, maxW*10))
	}
	return b.String()
}

func renderComment(c *archive.Comment, width int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Comment\n\n")
	fmt.Fprintf(&b, "User:     %s\n", c.Data.User.Username)
	fmt.Fprintf(&b, "Date:     %s\n", formatEpochMs(c.Data.Date))
	fmt.Fprintf(&b, "Resolved: %v\n", c.Data.Resolved)
	if c.Data.Text != "" {
		fmt.Fprintf(&b, "\n")
		maxW := max(width-4, 20)
		fmt.Fprintf(&b, "%s\n", truncate(c.Data.Text, maxW*10))
	}
	return b.String()
}

func formatEpochMs(s string) string {
	var ms int64
	if _, err := fmt.Sscanf(s, "%d", &ms); err != nil {
		return s
	}
	return time.UnixMilli(ms).Format("2006-01-02")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
