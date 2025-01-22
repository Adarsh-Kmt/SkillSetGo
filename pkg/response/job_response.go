package response

type JobOfferResponse struct {
	JobId       int32   `json:"job_id"`
	CompanyName string  `json:"company_name"`
	JobRole     string  `json:"job_role"`
	JobType     string  `json:"job_type"`
	CTC         float32 `json:"ctc"`
	SalaryTier  string  `json:"salary_tier"`
	Action      string  `json:"action"`
	ActionDate  string  `json:"action_date"`
	ActByDate   string  `json:"act_by_date"`
}
