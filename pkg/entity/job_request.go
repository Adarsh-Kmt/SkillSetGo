package entity

import (
	"time"

	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
)

var (
	validBranches = []string{"CSE", "ISE", "CY", "CD", "ECE", "EEE", "EIE", "ETE", "ME", "CV", "AS", "BT", "CH"}
)

type CreateJobRequest struct {
	JobRole          string   `json:"job_role"`
	Ctc              float32  `json:"ctc"`
	SalaryTier       string   `json:"salary_tier"`
	ApplyByDate      string   `json:"apply_by_date"`
	CgpaCutoff       float32  `json:"cgpa_cutoff"`
	EligibleBranches []string `json:"eligible_branches"`
	EligibleBatch    int      `json:"eligible_batch"`
	JobType          string   `json:"job_type"`
}

type PerformJobOfferActionRequest struct {
	JobId  int    `json:"job_id"`
	Action string `json:"action"`
}

type OfferJobRequest struct {
	StudentId int    `json:"student_id"`
	JobId     int    `json:"job_id"`
	ActByDate string `json:"act_by_date"`
}

func ValidateCreateJobRequest(request CreateJobRequest) (httpError *helper.HTTPError) {

	var (
		abd time.Time
		err error
	)

	errorMap := make(map[string]any, 0)

	if request.Ctc <= 0 {
		errorMap["ctc"] = "ctc cannot be negative/zero"
	}

	if request.SalaryTier == "" || (request.SalaryTier != "Open Dream" && request.SalaryTier != "Dream" && request.SalaryTier != "Mass Recruitment") {
		errorMap["salary_tier"] = "salary_tier cannot be empty, must be one of Open Dream, Dream, Mass Recruitment"
	}

	if (request.SalaryTier == "Open Dream" && request.Ctc <= 8) || (request.SalaryTier == "Dream" && request.Ctc <= 3) {
		errorMap["salary_tier"] = "salary tier does not match the CTC"
	}
	if request.ApplyByDate == "" {
		errorMap["apply_by_date"] = "apply_by_date cannot be empty"
	}

	if abd, err = time.Parse("2006-01-02 15:04:05", request.ApplyByDate); err != nil {
		errorMap["apply_by_date"] = "invalid apply_by_date format"
	}

	eligibleBranchesErrorList := make([]string, 0)
	if len(request.EligibleBranches) == 0 {
		eligibleBranchesErrorList = append(eligibleBranchesErrorList, "eligible_branches cannot be empty")
	}

	for _, branch := range request.EligibleBranches {
		if branch == "" {
			eligibleBranchesErrorList = append(eligibleBranchesErrorList, "branch cannot be empty")
			//return &helper.HTTPError{StatusCode: 400, Error: "branch cannot be empty"}
		}
		flag := false
		for _, validBranch := range validBranches {
			if branch == validBranch {
				flag = true
				break
			}
		}
		if !flag {
			eligibleBranchesErrorList = append(eligibleBranchesErrorList, "branch must be one of these valid branches")
			//return &helper.HTTPError{StatusCode: 400, Error: fmt.Sprintf("branch must be one of these valid branches %v", validBranches)}
		}

	}

	if len(eligibleBranchesErrorList) != 0 {
		errorMap["eligible_branches"] = eligibleBranchesErrorList
	}
	if request.CgpaCutoff < 0 || request.CgpaCutoff >= 10 {
		errorMap["cgpa_cutoff"] = "cgpa_cutoff cannot be negative/greater than 10"
		//return &helper.HTTPError{StatusCode: 400, Error: }
	}

	if errorMap["apply_by_date"] != nil {
		return &helper.HTTPError{StatusCode: 400, Error: errorMap}
	}
	if time.Now().After(abd) {
		errorMap["apply_by_date"] = "apply by date expired"
	}

	if len(errorMap) != 0 {
		return &helper.HTTPError{StatusCode: 400, Error: errorMap}

	}
	return nil
}

func ValidatePerformJobOfferActionRequest(request PerformJobOfferActionRequest) *helper.HTTPError {

	errorMap := make(map[string]string)
	if request.Action != "REJECT" && request.Action != "ACCEPT" {
		errorMap["action"] = "invalid action"
	}

	if request.JobId <= 0 {
		errorMap["job_id"] = "invalid job_id"
	}

	if len(errorMap) != 0 {
		return &helper.HTTPError{StatusCode: 400, Error: errorMap}
	}

	return nil
}

func ValidateOfferJobRequest(request OfferJobRequest) *helper.HTTPError {

	errorMap := make(map[string]string)

	if request.JobId <= 0 {
		errorMap["job_id"] = "invalid job_id"
	}
	if request.StudentId <= 0 {
		errorMap["student_id"] = "invalid student_id"
	}

	abd, err := time.Parse("2006-01-02 15:04:05", request.ActByDate)

	if err != nil {
		errorMap["act_by_date"] = "invalid act_by_date format"

		return &helper.HTTPError{StatusCode: 400, Error: errorMap}
	}

	if time.Now().After(abd) {
		errorMap["act_by_date"] = "apply by date expired"
	}
	if len(errorMap) != 0 {
		return &helper.HTTPError{StatusCode: 400, Error: errorMap}
	}

	return nil

}
