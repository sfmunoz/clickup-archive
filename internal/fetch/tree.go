package fetch

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type FetchTree struct {
	client *Client
}

func NewFetchTree(token string) *FetchTree {
	client := NewClient(token)
	return &FetchTree{client: client}
}

func (f *FetchTree) dumpTask(task api.Task, baseDir string) error {
	log.Info("Task", "id", task.ID, "name", task.Name)
	dir := filepath.Join(baseDir, task.ID)
	if err := jsonDump(task, dir); err != nil {
		return fmt.Errorf("dump task %s: %w", task.ID, err)
	}
	for _, sub := range task.Subtasks {
		if err := f.dumpTask(sub, baseDir); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) getTasks(listID, baseDir string) error {
	for page := 0; ; page++ {
		var resp api.TasksResponse
		path := fmt.Sprintf("/list/%s/task?include_closed=true&subtasks=true&page=%d", listID, page)
		if err := f.client.HttpGet(path, &resp); err != nil {
			return fmt.Errorf("fetch tasks page %d: %w", page, err)
		}
		for _, task := range resp.Tasks {
			if err := f.dumpTask(task, baseDir); err != nil {
				return err
			}
		}
		if len(resp.Tasks) == 0 {
			break
		}
	}
	return nil
}

func (f *FetchTree) getLists(folderID, baseDir string) error {
	var resp api.ListsResponse
	if err := f.client.HttpGet("/folder/"+folderID+"/list", &resp); err != nil {
		return fmt.Errorf("fetch lists: %w", err)
	}
	for _, list := range resp.Lists {
		log.Info("List", "id", list.ID, "name", list.Name)
		dir := filepath.Join(baseDir, list.ID)
		if err := jsonDump(list, dir); err != nil {
			return fmt.Errorf("dump list %s: %w", list.ID, err)
		}
		if err := f.getTasks(list.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) getFolders(spaceID, baseDir string) error {
	var resp api.FoldersResponse
	if err := f.client.HttpGet("/space/"+spaceID+"/folder", &resp); err != nil {
		return fmt.Errorf("fetch folders: %w", err)
	}
	for _, folder := range resp.Folders {
		log.Info("Folder", "id", folder.ID, "name", folder.Name)
		dir := filepath.Join(baseDir, folder.ID)
		if err := jsonDump(folder, dir); err != nil {
			return fmt.Errorf("dump folder %s: %w", folder.ID, err)
		}
		if err := f.getLists(folder.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) getSpaces(workspaceID, baseDir string) error {
	var resp api.SpacesResponse
	if err := f.client.HttpGet("/team/"+workspaceID+"/space", &resp); err != nil {
		return fmt.Errorf("fetch spaces: %w", err)
	}
	for _, space := range resp.Spaces {
		log.Info("Space", "id", space.ID, "name", space.Name)
		dir := filepath.Join(baseDir, space.ID)
		if err := jsonDump(space, dir); err != nil {
			return fmt.Errorf("dump space %s: %w", space.ID, err)
		}
		if err := f.getFolders(space.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchTree) Run(baseDir string) error {
	var resp api.WorkspacesResponse
	if err := f.client.HttpGet("/team", &resp); err != nil {
		return fmt.Errorf("fetch workspaces: %w", err)
	}
	for _, workspace := range resp.Workspaces {
		log.Info("Workspace", "id", workspace.ID, "name", workspace.Name)
		dir := filepath.Join(baseDir, workspace.ID)
		if err := jsonDump(workspace, dir); err != nil {
			return fmt.Errorf("dump workspace %s: %w", workspace.ID, err)
		}
		if err := f.getSpaces(workspace.ID, dir); err != nil {
			return err
		}
	}
	return nil
}
