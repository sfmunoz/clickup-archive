package tui

import (
	"testing"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

func TestItemsInitialStateShowsTwoLevels(t *testing.T) {
	items := newTestItems(t)

	assertVisibleTitles(t, items, []string{"Workspace", "Space"})

	workspace := itemAt(t, items, 0)
	if !workspace.IsExpandable() || !workspace.IsExpanded() {
		t.Fatalf("workspace should start expanded")
	}
	space := itemAt(t, items, 1)
	if !space.IsExpandable() || space.IsExpanded() {
		t.Fatalf("space should start collapsed")
	}
}

func TestItemsExpandSelectedSubtreeOneLevelAtATime(t *testing.T) {
	items := newTestItems(t)
	items.list.Select(1)

	items.expandSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder"})

	items.expandSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder", "List"})

	items.expandSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder", "List", "task-1"})
}

func TestItemsCollapseSelectedSubtreeOneLevelAtATime(t *testing.T) {
	items := newTestItems(t)
	items.list.Select(1)
	items.expandSelected()
	items.expandSelected()
	items.expandSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder", "List", "task-1"})

	items.collapseSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder", "List"})

	items.collapseSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder"})

	items.collapseSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space"})
}

func TestItemsIgnoreLeafExpandCollapse(t *testing.T) {
	items := newTestItems(t)
	items.list.Select(1)
	items.expandSelected()
	items.expandSelected()
	items.expandSelected()
	items.expandSelected()
	items.list.Select(5)

	items.expandSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder", "List", "task-1", "comment-1"})

	items.collapseSelected()
	assertVisibleTitles(t, items, []string{"Workspace", "Space", "Folder", "List", "task-1", "comment-1"})
}

func TestItemsPreserveSelectedNode(t *testing.T) {
	items := newTestItems(t)
	items.list.Select(1)
	selected := items.SelectedNode()

	items.expandSelected()

	if got := items.SelectedNode(); got != selected {
		t.Fatalf("selected node changed after expand")
	}
}

func newTestItems(t *testing.T) *Items {
	t.Helper()
	items, err := NewItems(testArchive())
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func testArchive() *archive.Archive {
	a := &archive.Archive{
		Data:     &archive.ArchiveData{Dir: "/tmp/clickup-archive-test"},
		Children: make([]*archive.Workspace, 0),
	}
	w := &archive.Workspace{Parent: a, Data: &api.Workspace{ID: "workspace-1", Name: "Workspace"}}
	s := &archive.Space{Parent: w, Data: &api.Space{ID: "space-1", Name: "Space"}}
	f := &archive.Folder{Parent: s, Data: &api.Folder{ID: "folder-1", Name: "Folder"}}
	l := &archive.List{Parent: f, Data: &api.List{ID: "list-1", Name: "List"}}
	task := &archive.Task{Parent: l, Data: &api.Task{ID: "task-1", Name: "Task"}}
	comment := &archive.Comment{Parent: task, Data: &api.Comment{ID: "comment-1", Text: "Comment"}}
	task.Children = []*archive.Comment{comment}
	l.Children = []*archive.Task{task}
	f.Children = []*archive.List{l}
	s.Children = []*archive.Folder{f}
	w.Children = []*archive.Space{s}
	a.Children = []*archive.Workspace{w}
	return a
}

func itemAt(t *testing.T, items *Items, idx int) *item {
	t.Helper()
	got, ok := items.list.Items()[idx].(*item)
	if !ok {
		t.Fatalf("item %d has unexpected type %T", idx, items.list.Items()[idx])
	}
	return got
}

func assertVisibleTitles(t *testing.T, items *Items, want []string) {
	t.Helper()
	gotItems := items.list.Items()
	if len(gotItems) != len(want) {
		t.Fatalf("visible item count = %d, want %d; titles = %v", len(gotItems), len(want), visibleTitles(items))
	}
	for idx, title := range want {
		if got := itemAt(t, items, idx).Title(); got != title {
			t.Fatalf("title %d = %q, want %q; titles = %v", idx, got, title, visibleTitles(items))
		}
	}
}

func visibleTitles(items *Items) []string {
	titles := make([]string, 0, len(items.list.Items()))
	for _, listItem := range items.list.Items() {
		if item, ok := listItem.(*item); ok {
			titles = append(titles, item.Title())
		}
	}
	return titles
}
