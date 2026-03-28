package api

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
