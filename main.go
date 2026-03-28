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
	for _, line := range strings.Split(string(data), "\n") {
		log.Info(line)
	}
	return nil
}

func main() {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		log.Fatal("CLICKUP_TOKEN env var is required")
	}
	var workspaces api.WorkspacesResponse
	if err := get(token, "/team", &workspaces); err != nil {
		log.Fatal("Failed to fetch workspaces", "err", err)
	}
	for _, workspace := range workspaces.Workspaces {
		log.Info("Workspace", "name", workspace.Name, "id", workspace.ID)
		jsonDump(workspace)
		var spaces api.SpacesResponse
		if err := get(token, "/team/"+workspace.ID+"/space", &spaces); err != nil {
			log.Fatal("Failed to fetch spaces", "err", err)
		}
		for _, space := range spaces.Spaces {
			log.Info("Space", "name", space.Name, "id", space.ID)
			jsonDump(space)
		}
	}
}
