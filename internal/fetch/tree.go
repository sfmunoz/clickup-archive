package fetch

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

type FetchTree struct {
	archive *archive.Archive
	update  bool
	client  *Client
}

func NewFetchTree(a *archive.Archive, update bool) (*FetchTree, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &FetchTree{
		archive: a,
		update:  update,
		client:  client,
	}, nil
}

func (f *FetchTree) dumpTask(task api.Task, archLi *archive.List) error {
	log.Info("Task", "id", task.ID, "name", task.Name)
	if _, err := archive.SaveTask(archLi, &task, f.update); err != nil {
		return fmt.Errorf("dump task %s: %w", task.ID, err)
	}
	for _, sub := range task.Subtasks {
		if err := f.dumpTask(sub, archLi); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) getTasks(listID string, archLi *archive.List) error {
	for page := 0; ; page++ {
		var resp api.TasksResponse
		path := fmt.Sprintf("/list/%s/task?include_closed=true&subtasks=true&page=%d", listID, page)
		if err := f.client.HttpGet(path, &resp); err != nil {
			return fmt.Errorf("fetch tasks page %d: %w", page, err)
		}
		for _, task := range resp.Tasks {
			if err := f.dumpTask(task, archLi); err != nil {
				return err
			}
		}
		if len(resp.Tasks) == 0 {
			break
		}
	}
	return nil
}

func (f *FetchTree) getLists(folderID string, archFo *archive.Folder) error {
	var resp api.ListsResponse
	if err := f.client.HttpGet("/folder/"+folderID+"/list", &resp); err != nil {
		return fmt.Errorf("fetch lists: %w", err)
	}
	for _, list := range resp.Lists {
		log.Info("List", "id", list.ID, "name", list.Name)
		archLi, err := archive.SaveList(archFo, &list, f.update)
		if err != nil {
			return err
		}
		if err := f.getTasks(list.ID, archLi); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) getFolders(spaceID string, archSp *archive.Space) error {
	var resp api.FoldersResponse
	if err := f.client.HttpGet("/space/"+spaceID+"/folder", &resp); err != nil {
		return fmt.Errorf("fetch folders: %w", err)
	}
	for _, folder := range resp.Folders {
		log.Info("Folder", "id", folder.ID, "name", folder.Name)
		archFo, err := archive.SaveFolder(archSp, &folder, f.update)
		if err != nil {
			return err
		}
		if err := f.getLists(folder.ID, archFo); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) getSpaces(workspaceID string, archWs *archive.Workspace) error {
	var resp api.SpacesResponse
	if err := f.client.HttpGet("/team/"+workspaceID+"/space", &resp); err != nil {
		return fmt.Errorf("fetch spaces: %w", err)
	}
	for _, space := range resp.Spaces {
		log.Info("Space", "id", space.ID, "name", space.Name)
		archSp, err := archive.SaveSpace(archWs, &space, f.update)
		if err != nil {
			return err
		}
		if err := f.getFolders(space.ID, archSp); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) Run() error {
	var resp api.WorkspacesResponse
	if err := f.client.HttpGet("/team", &resp); err != nil {
		return fmt.Errorf("fetch workspaces: %w", err)
	}
	for _, workspace := range resp.Workspaces {
		log.Info("Workspace", "id", workspace.ID, "name", workspace.Name)
		archWs, err := archive.SaveWorkspace(f.archive, &workspace, f.update)
		if err != nil {
			return err
		}
		if err := f.getSpaces(workspace.ID, archWs); err != nil {
			return err
		}
	}
	return nil
}
