package entity

import (
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"time"
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

	if request.SalaryTier == "" || (request.SalaryTier != "Open Dream" && request.SalaryTier != "Dream") {
		return &util.HTTPError{StatusCode: 400, Error: "salary_tier cannot be empty, must be one of Open Dream, Dream"}
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

	if request.CgpaCutoff < 0 || request.CgpaCutoff >= 10 {
		return &util.HTTPError{StatusCode: 400, Error: "cgpa_cutoff cannot be negative/greater than 10"}
	}
	return nil
}
