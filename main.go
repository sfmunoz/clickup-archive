package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
		fmt.Fprintln(os.Stderr, "CLICKUP_TOKEN env var is required")
		os.Exit(1)
	}

	var teams TeamsResponse
	if err := get(token, "/team", &teams); err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch workspaces: %v\n", err)
		os.Exit(1)
	}

	for _, team := range teams.Teams {
		fmt.Printf("Workspace: %s (id=%s)\n", team.Name, team.ID)
		var spaces SpacesResponse
		if err := get(token, "/team/"+team.ID+"/space", &spaces); err != nil {
			fmt.Fprintf(os.Stderr, "  failed to fetch spaces: %v\n", err)
			continue
		}
		for _, space := range spaces.Spaces {
			fmt.Printf("  Space: %s (id=%s)\n", space.Name, space.ID)
		}
	}
}
