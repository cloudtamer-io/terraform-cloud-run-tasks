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

	"github.com/cloudtamer-io/terraform-cloud-run-tasks/lib"
)

func main() {
	// If the port is passed as an argument, use it.
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fmt.Println("Webhook running on port:", port)
	http.HandleFunc("/", index)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/favicon.ico") {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Webhook called:", r.URL.Path)

	action := strings.TrimPrefix(r.URL.Path, "/")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	payload := new(lib.TFTaskRequest)

	err = json.Unmarshal(b, payload)
	if err != nil {
		http.NotFound(w, r)
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
	vars := new(lib.TFVarsResponse)
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

	// Change the token so it can send a status to the callback URL.
	rc2.Token = payload.AccessToken

	if len(projectID) == 0 {
		log.Println("terraform cloud env variable missing: CLOUDTAMERIO_PROJECT")

		tfresult := new(lib.TFResultRequest)
		tfresult.Data.Attributes.URL = "https://cloudtamer.zendesk.com/hc/en-us/articles/4408728893325"
		tfresult.Data.Attributes.Status = "failed"
		tfresult.Data.Attributes.Message = "You must set the environment variable, CLOUDTAMERIO_PROJECT, as a workspace variable."
		payback := new(lib.TFTaskResponse)
		err = tf.PATCH(payload.TaskResultCallbackURL, tfresult, payback)
		if err != nil {
			log.Printf("error on sending result: %v\n", err)
			return
		}

		return
	}

	fmt.Println("Project ID:", projectID)

	switch action {
	case "savings":
		sendSavings(w, ct, tf, payload, cturl, projectID)
	case "compliance":
		sendCompliance(w, ct, tf, payload, cturl, projectID)
	default:
		fmt.Println("No action matching:", action)
	}
}

func sendSavings(w http.ResponseWriter, ct *lib.CTClient, tf *lib.TerraformCloudClient, payload *lib.TFTaskRequest, cturl string, projectID string) {
	spend := new(lib.ProjectSpendReponse)
	err := ct.GET(fmt.Sprintf("/v3/project/%v/spend/monthly", projectID), spend)
	if err != nil {
		log.Printf("error on getting spend: %v\n", err)
		return
	}

	fmt.Println("Spend:", spend.Data.Spend)
	fmt.Println("Estimate:", spend.Data.Estimate)

	savings := new(lib.CostSavingsResponse)
	err = ct.GET(fmt.Sprintf("/v1/cost-savings/cost-and-savings?csp_type_id=0&service_id=0&project_id=%v&include_dismissed=false&projectFilter=true", projectID), savings)
	if err != nil {
		log.Printf("error on getting savings: %v\n", err)
		return
	}

	totalSavings := savings.Data.CurrentMonthlyCost - savings.Data.PotentialMonthlyCost

	fmt.Println("Potential savings:", totalSavings)

	// Send response back to Terraform Cloud.
	tfresult := new(lib.TFResultRequest)
	tfresult.Data.Attributes.URL = fmt.Sprintf("%v/portal/project/%v/savings-opportunity", cturl, projectID)
	tfresult.Data.Attributes.Status = "passed"
	tfresult.Data.Attributes.Message = fmt.Sprintf("For project, Nimbus, the monthly forecast is $%v. You could be saving $%v per month through savings opportunities. Click Details to view them.", math.Round(savings.Data.CurrentMonthlyCost), math.Round(totalSavings))
	payback := new(lib.TFTaskResponse)
	err = tf.PATCH(payload.TaskResultCallbackURL, tfresult, payback)
	if err != nil {
		log.Printf("error on sending result: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Success: %v", payback.Data.Attributes.Message)
}

func sendCompliance(w http.ResponseWriter, ct *lib.CTClient, tf *lib.TerraformCloudClient, payload *lib.TFTaskRequest, cturl string, projectID string) {
	compliance := new(lib.ComplianceResponse)
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

	tfresult := new(lib.TFResultRequest)
	tfresult.Data.Attributes.URL = fmt.Sprintf("%v/portal/project/%v/compliance", cturl, projectID)
	tfresult.Data.Attributes.Status = "passed"

	tfresult.Data.Attributes.Message = fmt.Sprintf("For project, Nimbus, there are %v compliance findings. | Critical: %v | High: %v | Medium: %v | Low: %v | Info: %v | Click Details to view them.", len(compliance.Data.Items), critical, high, medium, low, info)

	if critical > 0 {
		tfresult.Data.Attributes.Status = "failed"
		tfresult.Data.Attributes.Message += " You cannot have any critical findings."
	}

	payback := new(lib.TFTaskResponse)

	err = tf.PATCH(payload.TaskResultCallbackURL, tfresult, payback)
	if err != nil {
		log.Printf("error on sending result: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Success: %v", payback.Data.Attributes.Message)
}
