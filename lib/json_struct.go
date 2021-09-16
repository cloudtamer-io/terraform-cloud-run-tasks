package lib

import "time"

// ProjectSpendReponse returns the monthly spend for a project.
type ProjectSpendReponse struct {
	Status int `json:"status"`
	Data   struct {
		Spend    float64 `json:"spend"`
		Estimate float64 `json:"estimate"`
	} `json:"data"`
}

// TFResultRequest sends the result.
type TFResultRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Status  string `json:"status"`
			URL     string `json:"url"`
			Message string `json:"message"`
		} `json:"attributes"`
	} `json:"data"`
}

// TFTaskRequest sends the payload when a plan runs.
type TFTaskRequest struct {
	PayloadVersion             int       `json:"payload_version"`
	AccessToken                string    `json:"access_token"`
	TaskResultID               string    `json:"task_result_id"`
	TaskResultEnforcementLevel string    `json:"task_result_enforcement_level"`
	TaskResultCallbackURL      string    `json:"task_result_callback_url"`
	RunAppURL                  string    `json:"run_app_url"`
	RunID                      string    `json:"run_id"`
	RunMessage                 string    `json:"run_message"`
	RunCreatedAt               time.Time `json:"run_created_at"`
	RunCreatedBy               string    `json:"run_created_by"`
	WorkspaceID                string    `json:"workspace_id"`
	WorkspaceName              string    `json:"workspace_name"`
	WorkspaceAppURL            string    `json:"workspace_app_url"`
	OrganizationName           string    `json:"organization_name"`
	PlanJSONAPIURL             string    `json:"plan_json_api_url"`
	VcsRepoURL                 string    `json:"vcs_repo_url"`
	VcsBranch                  string    `json:"vcs_branch"`
	VcsPullRequestURL          string    `json:"vcs_pull_request_url"`
	VcsCommitURL               string    `json:"vcs_commit_url"`
}

// TFTaskResponse is the response back from the callback.
type TFTaskResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Message          string `json:"message"`
			Status           string `json:"status"`
			StatusTimestamps struct {
				RunningAt time.Time `json:"running-at"`
				PassedAt  time.Time `json:"passed-at"`
			} `json:"status-timestamps"`
			URL                  string    `json:"url"`
			CreatedAt            time.Time `json:"created-at"`
			UpdatedAt            time.Time `json:"updated-at"`
			EventHookID          string    `json:"event-hook-id"`
			EventHookName        string    `json:"event-hook-name"`
			EventHookURL         string    `json:"event-hook-url"`
			TaskID               string    `json:"task-id"`
			TaskEnforcementLevel string    `json:"task-enforcement-level"`
		} `json:"attributes"`
		Relationships struct {
			PreApply struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"pre-apply"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

// TFVarsResponse returns the environment variables for the workspace.
type TFVarsResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Key         string      `json:"key"`
			Value       interface{} `json:"value"`
			Sensitive   bool        `json:"sensitive"`
			Category    string      `json:"category"`
			Hcl         bool        `json:"hcl"`
			CreatedAt   time.Time   `json:"created-at"`
			Description interface{} `json:"description"`
		} `json:"attributes"`
		Relationships struct {
			Configurable struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"configurable"`
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

// CostSavingsResponse returns the cost savings for a project.
type CostSavingsResponse struct {
	Status int `json:"status"`
	Data   struct {
		CurrentMonthlyCost   float64 `json:"current_monthly_cost"`
		PotentialMonthlyCost float64 `json:"potential_monthly_cost"`
		DecommissionSavings  float64 `json:"decommission_savings"`
		RightsizingSavings   float64 `json:"rightsizing_savings"`
		MonthCount           int     `json:"month_count"`
	} `json:"data"`
}

// ComplianceResponse returns the compliance for a project.
type ComplianceResponse struct {
	Status int `json:"status"`
	Data   struct {
		Total int `json:"total"`
		Items []struct {
			Finding struct {
				ID                    int    `json:"id"`
				HashIdentifier        int    `json:"hash_identifier"`
				ComplianceCheckScanID int    `json:"compliance_check_scan_id"`
				ResourceType          string `json:"resource_type"`
				ResourceName          string `json:"resource_name"`
			} `json:"finding"`
			StandardID     int    `json:"standard_id"`
			StandardName   string `json:"standard_name"`
			CheckID        int    `json:"check_id"`
			CheckName      string `json:"check_name"`
			SeverityTypeID int    `json:"severity_type_id"`
			AccountID      int    `json:"account_id"`
			AccountNumber  string `json:"account_number"`
			AccountName    string `json:"account_name"`
			ProjectID      int    `json:"project_id"`
			ProjectName    string `json:"project_name"`
			ParentOuID     int    `json:"parent_ou_id"`
			ParentOuName   string `json:"parent_ou_name"`
			Region         string `json:"region"`
			CreatedAt      string `json:"created_at"`
			ArchivedAt     string `json:"archived_at"`
		} `json:"items"`
	} `json:"data"`
}
