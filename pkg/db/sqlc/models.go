// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CompanyTable struct {
	CompanyID   int32  `json:"company_id"`
	CompanyName string `json:"company_name"`
	PocName     string `json:"poc_name"`
	PocPhno     string `json:"poc_phno"`
	Industry    string `json:"industry"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type JobTable struct {
	JobID            int32            `json:"job_id"`
	CompanyID        int32            `json:"company_id"`
	JobRole          string           `json:"job_role"`
	Ctc              float32          `json:"ctc"`
	SalaryTier       string           `json:"salary_tier"`
	ApplyByDate      pgtype.Timestamp `json:"apply_by_date"`
	CgpaCutoff       float32          `json:"cgpa_cutoff"`
	EligibleBranches []string         `json:"eligible_branches"`
}

type StudentJobApplicationTable struct {
	StudentID     int32            `json:"student_id"`
	JobID         int32            `json:"job_id"`
	AppliedOnDate pgtype.Timestamp `json:"applied_on_date"`
}

type StudentJobInterviewTable struct {
	StudentID      int32            `json:"student_id"`
	Venue          string           `json:"venue"`
	InterviewDate  pgtype.Timestamp `json:"interview_date"`
	InterviewRound int32            `json:"interview_round"`
	Result         string           `json:"result"`
}

type StudentOfferTable struct {
	StudentID  int32            `json:"student_id"`
	JobID      int32            `json:"job_id"`
	Action     string           `json:"action"`
	ActionDate pgtype.Timestamp `json:"action_date"`
	ActByDate  pgtype.Timestamp `json:"act_by_date"`
}

type StudentTable struct {
	StudentID         int32   `json:"student_id"`
	Usn               string  `json:"usn"`
	Name              string  `json:"name"`
	Branch            string  `json:"branch"`
	Cgpa              float32 `json:"cgpa"`
	NumActiveBacklogs int32   `json:"num_active_backlogs"`
	EmailID           string  `json:"email_id"`
	CounsellorEmailID string  `json:"counsellor_email_id"`
}
