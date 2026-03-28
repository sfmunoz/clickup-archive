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

func httpGetOnce(token, path string, out any) error {
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
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

func httpGet(token, path string, out any) error {
	var lastErr error
	for attempt := range httpGetRetries {
		if attempt > 0 {
			time.Sleep(httpGetRetryDelay)
		}
		if err := httpGetOnce(token, path, out); err != nil {
			lastErr = err
			log.Warn("httpGet failed, retrying", "attempt", attempt+1, "err", err)
			continue
		}
		return nil
	}
	return lastErr
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

func getLists(token, folderID, baseDir string) error {
	var resp api.ListsResponse
	if err := httpGet(token, "/folder/"+folderID+"/list", &resp); err != nil {
		return fmt.Errorf("fetch lists: %w", err)
	}
	for _, list := range resp.Lists {
		log.Info("List", "name", list.Name, "id", list.ID)
		dir := filepath.Join(baseDir, list.ID)
		if err := jsonDump(list, dir); err != nil {
			return fmt.Errorf("dump list %s: %w", list.ID, err)
		}
	}
	return nil
}

func getFolders(token, spaceID, baseDir string) error {
	var resp api.FoldersResponse
	if err := httpGet(token, "/space/"+spaceID+"/folder", &resp); err != nil {
		return fmt.Errorf("fetch folders: %w", err)
	}
	for _, folder := range resp.Folders {
		log.Info("Folder", "name", folder.Name, "id", folder.ID)
		dir := filepath.Join(baseDir, folder.ID)
		if err := jsonDump(folder, dir); err != nil {
			return fmt.Errorf("dump folder %s: %w", folder.ID, err)
		}
		if err := getLists(token, folder.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func getSpaces(token, workspaceID, baseDir string) error {
	var resp api.SpacesResponse
	if err := httpGet(token, "/team/"+workspaceID+"/space", &resp); err != nil {
		return fmt.Errorf("fetch spaces: %w", err)
	}
	for _, space := range resp.Spaces {
		log.Info("Space", "name", space.Name, "id", space.ID)
		dir := filepath.Join(baseDir, space.ID)
		if err := jsonDump(space, dir); err != nil {
			return fmt.Errorf("dump space %s: %w", space.ID, err)
		}
		if err := getFolders(token, space.ID, dir); err != nil {
			return err
		}
	}
	return nil
}

func getWorkspaces(token, baseDir string) error {
	var resp api.WorkspacesResponse
	if err := httpGet(token, "/team", &resp); err != nil {
		return fmt.Errorf("fetch workspaces: %w", err)
	}
	for _, workspace := range resp.Workspaces {
		log.Info("Workspace", "name", workspace.Name, "id", workspace.ID)
		dir := filepath.Join(baseDir, workspace.ID)
		if err := jsonDump(workspace, dir); err != nil {
			return fmt.Errorf("dump workspace %s: %w", workspace.ID, err)
		}
		if err := getSpaces(token, workspace.ID, dir); err != nil {
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
	if err := getWorkspaces(token, outputDir); err != nil {
		log.Fatal("Failed", "err", err)
	}
}
