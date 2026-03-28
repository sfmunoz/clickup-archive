package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const baseURL = "https://api.clickup.com/api/v2"

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type Space struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

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

func main() {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		log.Fatal("CLICKUP_TOKEN env var is required")
	}

	var teams TeamsResponse
	if err := get(token, "/team", &teams); err != nil {
		log.Fatal("Failed to fetch workspaces", "err", err)
	}

	for _, team := range teams.Teams {
		log.Info("Workspace", "name", team.Name, "id", team.ID)
		var spaces SpacesResponse
		if err := get(token, "/team/"+team.ID+"/space", &spaces); err != nil {
			log.Fatal("Failed to fetch spaces", "err", err)
		}
		for _, space := range spaces.Spaces {
			log.Info("Space", "name", space.Name, "id", space.ID)
		}
	}
}
