package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const outputDir = "output"
const baseURL = "https://api.clickup.com/api/v2"

func httpGet(token, path string, out any) error {
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}
	return json.Unmarshal(body, out)
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

func getLists(token, folderID, baseDir string) {
	var resp api.ListsResponse
	if err := httpGet(token, "/folder/"+folderID+"/list", &resp); err != nil {
		log.Fatal("Failed to fetch lists", "err", err)
	}
	for _, list := range resp.Lists {
		log.Info("List", "name", list.Name, "id", list.ID)
		dir := filepath.Join(baseDir, list.ID)
		jsonDump(list, dir)
	}
}

func getFolders(token, spaceID, baseDir string) {
	var resp api.FoldersResponse
	if err := httpGet(token, "/space/"+spaceID+"/folder", &resp); err != nil {
		log.Fatal("Failed to fetch folders", "err", err)
	}
	for _, folder := range resp.Folders {
		log.Info("Folder", "name", folder.Name, "id", folder.ID)
		dir := filepath.Join(baseDir, folder.ID)
		jsonDump(folder, dir)
		getLists(token, folder.ID, dir)
	}
}

func getSpaces(token, workspaceID, baseDir string) {
	var resp api.SpacesResponse
	if err := httpGet(token, "/team/"+workspaceID+"/space", &resp); err != nil {
		log.Fatal("Failed to fetch spaces", "err", err)
	}
	for _, space := range resp.Spaces {
		log.Info("Space", "name", space.Name, "id", space.ID)
		dir := filepath.Join(baseDir, space.ID)
		jsonDump(space, dir)
		getFolders(token, space.ID, dir)
	}
}

func getWorkspaces(token, baseDir string) {
	var resp api.WorkspacesResponse
	if err := httpGet(token, "/team", &resp); err != nil {
		log.Fatal("Failed to fetch workspaces", "err", err)
	}
	for _, workspace := range resp.Workspaces {
		log.Info("Workspace", "name", workspace.Name, "id", workspace.ID)
		dir := filepath.Join(baseDir, workspace.ID)
		jsonDump(workspace, dir)
		getSpaces(token, workspace.ID, dir)
	}
}

func main() {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		log.Fatal("CLICKUP_TOKEN env var is required")
	}
	getWorkspaces(token, outputDir)
}
