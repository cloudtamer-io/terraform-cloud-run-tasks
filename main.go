package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudtamer-io/terraform-cloud-run-tasks/lib"
)

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

func main() {
	fmt.Println("Webhook running.")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/favicon.ico") {
			http.NotFound(w, r)
			return
		}

		action := r.URL.Query().Get("action")
		fmt.Println("Webhook:", r.URL.Path, action)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			//fmt.Println("error:", err.Error())
			return
		}

		payload := new(TFTaskRequest)

		err = json.Unmarshal(b, payload)
		if err != nil {
			//fmt.Println("payload error:", err.Error())
			return
		}

		fmt.Printf("Payload: %v\n", string(b))
		fmt.Fprint(w, "OK\n")

		cturl := os.Getenv("CLOUDTAMERIO_URL")
		if len(cturl) == 0 {
			log.Println("env variable missing: CLOUDTAMERIO_URL")
			return
		}

		ctapi := os.Getenv("CLOUDTAMERIO_APIKEY")
		if len(ctapi) == 0 {
			log.Println("env variable missing: CLOUDTAMERIO_APIKEY")
			return
		}

		tfapi := os.Getenv("TERRAFORM_APIKEY")
		if len(tfapi) == 0 {
			log.Println("env variable missing: TERRAFORM_APIKEY")
			return
		}

		rc := lib.NewRequestClient(cturl+"/api", ctapi, true)
		ct := lib.NewCTClient(rc)

		rc2 := lib.NewRequestClient("https://app.terraform.io/api/", tfapi, true)
		rc2.ContentType = "application/vnd.api+json"
		tf := lib.NewTerraformCloudClient(rc2)

		// Pull the environment variables from Terraform.
		vars := new(TFVarsResponse)
		err = tf.GET(fmt.Sprintf("v2/workspaces/%v/vars", payload.WorkspaceID), vars)
		if err != nil {
			log.Printf("error on sending result: %v\n", err)
			return
		}

		projectID := ""
		for _, v := range vars.Data {
			if v.Attributes.Key == "CLOUDTAMERIO_PROJECT" {
				if v.Attributes.Value != nil {
					projectID = fmt.Sprint(v.Attributes.Value)
				}

				break
			}
		}

		if len(projectID) == 0 {
			log.Println("terraform cloud env variable missing: CLOUDTAMERIO_PROJECT")

			tfresult := new(TFResultRequest)
			tfresult.Data.Attributes.URL = ""
			tfresult.Data.Attributes.Status = "failed"
			tfresult.Data.Attributes.Message = "You must set the environment variable, CLOUDTAMERIO_PROJECT, as a workspace variable."
			payback := new(TFTaskResponse)
			err = tf.PATCH(payload.TaskResultCallbackURL, tfresult, payback)
			if err != nil {
				log.Printf("error on sending result: %v\n", err)
				return
			}

			return
		}

		fmt.Println("Project ID:", projectID)

		// Change the token so it can send a status to the callback URL.
		rc2.Token = payload.AccessToken

		switch action {
		case "savings":
			sendSavings(w, ct, tf, payload, cturl, projectID)
		case "compliance":
			sendCompliance(w, ct, tf, payload, cturl, projectID)
		}
	})

	log.Fatalln(http.ListenAndServe(":8080", mux))
}

func sendSavings(w http.ResponseWriter, ct *lib.CTClient, tf *lib.TerraformCloudClient, payload *TFTaskRequest, cturl string, projectID string) {
	spend := new(ProjectSpendReponse)
	err := ct.GET(fmt.Sprintf("/v3/project/%v/spend/monthly", projectID), spend)
	if err != nil {
		log.Printf("error on getting spend: %v\n", err)
		return
	}

	fmt.Println("Spend:", spend.Data.Spend)
	fmt.Println("Estimate:", spend.Data.Estimate)

	savings := new(CostSavingsResponse)
	err = ct.GET(fmt.Sprintf("/v1/cost-savings/cost-and-savings?csp_type_id=0&service_id=0&project_id=%v&include_dismissed=false&projectFilter=true", projectID), savings)
	if err != nil {
		log.Printf("error on getting savings: %v\n", err)
		return
	}

	totalSavings := savings.Data.CurrentMonthlyCost - savings.Data.PotentialMonthlyCost

	fmt.Println("Potential savings:", totalSavings)

	// Send response back to Terraform Cloud.
	tfresult := new(TFResultRequest)
	tfresult.Data.Attributes.URL = fmt.Sprintf("%v/portal/project/%v/savings-opportunity", cturl, projectID)
	tfresult.Data.Attributes.Status = "passed"
	tfresult.Data.Attributes.Message = fmt.Sprintf("For project, Nimbus, the monthly forecast is $%v. You could be saving $%v per month through savings opportunities. Click Details to view them.", math.Round(savings.Data.CurrentMonthlyCost), math.Round(totalSavings))
	payback := new(TFTaskResponse)
	err = tf.PATCH(payload.TaskResultCallbackURL, tfresult, payback)
	if err != nil {
		log.Printf("error on sending result: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Success: %v", payback.Data.Attributes.Message)
}

func sendCompliance(w http.ResponseWriter, ct *lib.CTClient, tf *lib.TerraformCloudClient, payload *TFTaskRequest, cturl string, projectID string) {
	compliance := new(ComplianceResponse)
	err := ct.GET(fmt.Sprintf("/v4/compliance/finding?project_id=%v&finding_type=active", projectID), compliance)
	if err != nil {
		log.Printf("error on getting savings: %v\n", err)
		return
	}

	fmt.Println("Compliance checks:", len(compliance.Data.Items))

	critical := 0
	high := 0
	medium := 0
	low := 0
	info := 0

	for _, v := range compliance.Data.Items {
		switch v.SeverityTypeID {
		case 1:
			info++
		case 2:
			low++
		case 3:
			medium++
		case 4:
			high++
		case 5:
			critical++
		}
	}

	tfresult := new(TFResultRequest)
	tfresult.Data.Attributes.URL = fmt.Sprintf("%v/portal/project/%v/compliance", cturl, projectID)
	tfresult.Data.Attributes.Status = "passed"

	tfresult.Data.Attributes.Message = fmt.Sprintf("For project, Nimbus, there are %v compliance findings. | Critical: %v | High: %v | Medium: %v | Low: %v | Info: %v | Click Details to view them.", len(compliance.Data.Items), critical, high, medium, low, info)

	if critical > 0 {
		tfresult.Data.Attributes.Status = "failed"
		tfresult.Data.Attributes.Message += " You cannot have any critical findings."
	}

	payback := new(TFTaskResponse)

	err = tf.PATCH(payload.TaskResultCallbackURL, tfresult, payback)
	if err != nil {
		log.Printf("error on sending result: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Success: %v", payback.Data.Attributes.Message)
}
