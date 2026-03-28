package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const baseURL = "https://api.clickup.com/api/v2"

func get(token, path string, out any) error {
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

func jsonDump(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		log.Info(line)
	}
	return nil
}

func getLists(token, folderID string) {
	var resp api.ListsResponse
	if err := get(token, "/folder/"+folderID+"/list", &resp); err != nil {
		log.Fatal("Failed to fetch lists", "err", err)
	}
	for _, list := range resp.Lists {
		log.Info("List", "name", list.Name, "id", list.ID)
		jsonDump(list)
	}
}

func getFolders(token, spaceID string) {
	var resp api.FoldersResponse
	if err := get(token, "/space/"+spaceID+"/folder", &resp); err != nil {
		log.Fatal("Failed to fetch folders", "err", err)
	}
	for _, folder := range resp.Folders {
		log.Info("Folder", "name", folder.Name, "id", folder.ID)
		jsonDump(folder)
		getLists(token, folder.ID)
	}
}

func getSpaces(token, workspaceID string) {
	var resp api.SpacesResponse
	if err := get(token, "/team/"+workspaceID+"/space", &resp); err != nil {
		log.Fatal("Failed to fetch spaces", "err", err)
	}
	for _, space := range resp.Spaces {
		log.Info("Space", "name", space.Name, "id", space.ID)
		jsonDump(space)
		getFolders(token, space.ID)
	}
}

func getWorkspaces(token string) {
	var resp api.WorkspacesResponse
	if err := get(token, "/team", &resp); err != nil {
		log.Fatal("Failed to fetch workspaces", "err", err)
	}
	for _, workspace := range resp.Workspaces {
		log.Info("Workspace", "name", workspace.Name, "id", workspace.ID)
		jsonDump(workspace)
		getSpaces(token, workspace.ID)
	}
}

func main() {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		log.Fatal("CLICKUP_TOKEN env var is required")
	}
	getWorkspaces(token)
}
