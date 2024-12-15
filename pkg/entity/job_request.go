package entity

import (
	"fmt"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"time"
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
}

func ValidateCreateJobRequest(request CreateJobRequest) (httpError *util.HTTPError) {

	var (
		abd time.Time
		err error
	)
	if request.Ctc <= 0 {
		return &util.HTTPError{StatusCode: 400, Error: "ctc cannot be negative/zero"}
	}

	if request.SalaryTier == "" || (request.SalaryTier != "Open Dream" && request.SalaryTier != "Dream" && request.SalaryTier != "Mass Recruitment") {
		return &util.HTTPError{StatusCode: 400, Error: "salary_tier cannot be empty, must be one of Open Dream, Dream"}
	}

	if (request.SalaryTier == "Open Dream" && request.Ctc <= 8) || (request.SalaryTier == "Dream" && request.Ctc <= 3) {
		return &util.HTTPError{StatusCode: 400, Error: "salary tier does not match the CTC"}
	}
	if request.ApplyByDate == "" {
		return &util.HTTPError{StatusCode: 400, Error: "apply_by_date cannot be empty"}
	}

	if abd, err = time.Parse("2006-01-02 15:04:05", request.ApplyByDate); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "invalid apply_by_date format"}
	}

	if time.Now().After(abd) {
		return &util.HTTPError{StatusCode: 400, Error: "apply by date must be in the future"}
	}

	if len(request.EligibleBranches) == 0 {
		return &util.HTTPError{StatusCode: 400, Error: "eligible_branches cannot be empty"}
	}

	for _, branch := range request.EligibleBranches {
		if branch == "" {
			return &util.HTTPError{StatusCode: 400, Error: "branch cannot be empty"}
		}
		flag := false
		for _, validBranch := range validBranches {
			if branch == validBranch {
				flag = true
				break
			}
		}
		if flag == false {
			return &util.HTTPError{StatusCode: 400, Error: fmt.Sprintf("branch must be one of these valid branches %v", validBranches)}
		}

	}

	if request.CgpaCutoff < 0 || request.CgpaCutoff >= 10 {
		return &util.HTTPError{StatusCode: 400, Error: "cgpa_cutoff cannot be negative/greater than 10"}
	}
	return nil
}
