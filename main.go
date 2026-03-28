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

type TeamMemberUser struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Color             string `json:"color"`
	ProfilePicture    string `json:"profilePicture"`
	Initials          string `json:"initials"`
	WeekStartDay      int    `json:"week_start_day"`
	GlobalFontSupport bool   `json:"global_font_support"`
	Timezone          string `json:"timezone"`
}

type TeamMember struct {
	User TeamMemberUser `json:"user"`
}

type Team struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Color   string       `json:"color"`
	Avatar  string       `json:"avatar"`
	Members []TeamMember `json:"members"`
}

type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type SpaceStatus struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Orderindex int    `json:"orderindex"`
	Color      string `json:"color"`
}

type SpaceFeatureDueDates struct {
	Enabled            bool `json:"enabled"`
	StartDate          bool `json:"start_date"`
	RemapDueDates      bool `json:"remap_due_dates"`
	RemapClosedDueDate bool `json:"remap_closed_due_date"`
}

type SpaceFeatureTimeTracking struct {
	Enabled bool `json:"enabled"`
}

type SpaceFeatureTimeEstimates struct {
	Enabled     bool `json:"enabled"`
	Rollup      bool `json:"rollup"`
	PerAssignee bool `json:"per_assignee"`
}

type SpaceFeaturePriority struct {
	ID         string `json:"id"`
	Priority   string `json:"priority"`
	Color      string `json:"color"`
	Orderindex string `json:"orderindex"`
}

type SpaceFeaturePriorities struct {
	Enabled    bool                   `json:"enabled"`
	Priorities []SpaceFeaturePriority `json:"priorities"`
}

type SpaceFeatureEnabled struct {
	Enabled bool `json:"enabled"`
}

type SpaceFeatures struct {
	DueDates                   SpaceFeatureDueDates      `json:"due_dates"`
	Sprints                    SpaceFeatureEnabled       `json:"sprints"`
	Points                     SpaceFeatureEnabled       `json:"points"`
	CustomItems                SpaceFeatureEnabled       `json:"custom_items"`
	Priorities                 SpaceFeaturePriorities    `json:"priorities"`
	Tags                       SpaceFeatureEnabled       `json:"tags"`
	CheckUnresolvedBeforeClose SpaceFeatureEnabled       `json:"check_unresolved_before_close"`
	Zoom                       SpaceFeatureEnabled       `json:"zoom"`
	Milestones                 SpaceFeatureEnabled       `json:"milestones"`
	CustomFields               SpaceFeatureEnabled       `json:"custom_fields"`
	RemapDependencies          SpaceFeatureEnabled       `json:"remap_dependencies"`
	DependencyWarning          SpaceFeatureEnabled       `json:"dependency_warning"`
	MultipleAssignees          SpaceFeatureEnabled       `json:"multiple_assignees"`
	Emails                     SpaceFeatureEnabled       `json:"emails"`
	TimeTracking               SpaceFeatureTimeTracking  `json:"time_tracking"`
	TimeEstimates              SpaceFeatureTimeEstimates `json:"time_estimates"`
	Checklists                 SpaceFeatureEnabled       `json:"checklists"`
	Portfolios                 SpaceFeatureEnabled       `json:"portfolios"`
}

type Space struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Private           bool          `json:"private"`
	Statuses          []SpaceStatus `json:"statuses"`
	MultipleAssignees bool          `json:"multiple_assignees"`
	Features          SpaceFeatures `json:"features"`
	Archived          bool          `json:"archived"`
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

func jsonDump(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
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
		jsonDump(team)
		var spaces SpacesResponse
		if err := get(token, "/team/"+team.ID+"/space", &spaces); err != nil {
			log.Fatal("Failed to fetch spaces", "err", err)
		}
		for _, space := range spaces.Spaces {
			log.Info("Space", "name", space.Name, "id", space.ID)
			jsonDump(space)
		}
	}
}
