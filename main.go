package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const outputDir = "output"
const baseURL = "https://api.clickup.com/api/v2"
const httpGetRetries = 5
const httpGetRetryDelay = time.Second

type Client struct {
	token string
}

func (c *Client) httpGetOnce(path string, out any) error {
	time.Sleep(650 * time.Millisecond) // limit = 100 request/minute → 0.6 sec/request
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}
	return json.Unmarshal(body, out)
}

func (c *Client) httpGet(path string, out any) error {
	for attempt := 1; attempt <= httpGetRetries; attempt++ {
		err := c.httpGetOnce(path, out)
		if err == nil {
			break
		}
		if attempt == httpGetRetries {
			return err
		}
		log.Warn("httpGet failed, retrying", "attempt", attempt, "err", err)
		time.Sleep(httpGetRetryDelay)
	}
	return nil
}

func jsonDump(v any, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		log.Debug(line)
	}
	return os.WriteFile(filepath.Join(dir, "index.json"), data, 0o644)
}

func (c *Client) dumpTask(task api.Task, baseDir string) error {
	log.Info("Task", "name", task.Name, "id", task.ID)
	dir := filepath.Join(baseDir, task.ID)
	if err := jsonDump(task, dir); err != nil {
		return fmt.Errorf("dump task %s: %w", task.ID, err)
	}
	for _, sub := range task.Subtasks {
		if err := c.dumpTask(sub, baseDir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getTasks(listID, baseDir string) error {
	var resp api.TasksResponse
	if err := c.httpGet("/list/"+listID+"/task?include_closed=true&subtasks=true", &resp); err != nil {
		return fmt.Errorf("fetch tasks: %w", err)
	}
	for _, task := range resp.Tasks {
		if err := c.dumpTask(task, baseDir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getLists(folderID, baseDir string) error {
	var resp api.ListsResponse
	if err := c.httpGet("/folder/"+folderID+"/list", &resp); err != nil {
		return fmt.Errorf("fetch lists: %w", err)
	}
	for _, list := range resp.Lists {
		log.Info("List", "name", list.Name, "id", list.ID)
		dir := filepath.Join(baseDir, list.ID)
		if err := jsonDump(list, dir); err != nil {
			return fmt.Errorf("dump list %s: %w", list.ID, err)
		}
		if err := c.getTasks(list.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getFolders(spaceID, baseDir string) error {
	var resp api.FoldersResponse
	if err := c.httpGet("/space/"+spaceID+"/folder", &resp); err != nil {
		return fmt.Errorf("fetch folders: %w", err)
	}
	for _, folder := range resp.Folders {
		log.Info("Folder", "name", folder.Name, "id", folder.ID)
		dir := filepath.Join(baseDir, folder.ID)
		if err := jsonDump(folder, dir); err != nil {
			return fmt.Errorf("dump folder %s: %w", folder.ID, err)
		}
		if err := c.getLists(folder.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getSpaces(workspaceID, baseDir string) error {
	var resp api.SpacesResponse
	if err := c.httpGet("/team/"+workspaceID+"/space", &resp); err != nil {
		return fmt.Errorf("fetch spaces: %w", err)
	}
	for _, space := range resp.Spaces {
		log.Info("Space", "name", space.Name, "id", space.ID)
		dir := filepath.Join(baseDir, space.ID)
		if err := jsonDump(space, dir); err != nil {
			return fmt.Errorf("dump space %s: %w", space.ID, err)
		}
		if err := c.getFolders(space.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getWorkspaces(baseDir string) error {
	var resp api.WorkspacesResponse
	if err := c.httpGet("/team", &resp); err != nil {
		return fmt.Errorf("fetch workspaces: %w", err)
	}
	for _, workspace := range resp.Workspaces {
		log.Info("Workspace", "name", workspace.Name, "id", workspace.ID)
		dir := filepath.Join(baseDir, workspace.ID)
		if err := jsonDump(workspace, dir); err != nil {
			return fmt.Errorf("dump workspace %s: %w", workspace.ID, err)
		}
		if err := c.getSpaces(workspace.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		log.Fatal("CLICKUP_TOKEN env var is required")
	}
	c := &Client{token: token}
	if err := c.getWorkspaces(outputDir); err != nil {
		log.Fatal("Failed", "err", err)
	}
}
