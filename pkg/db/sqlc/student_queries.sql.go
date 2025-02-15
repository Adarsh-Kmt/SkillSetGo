// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: student_queries.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const applyForJob = `-- name: ApplyForJob :exec
INSERT INTO student_job_application_table(student_id, job_id, applied_on_date)
VALUES($1, $2, NOW())
`

type ApplyForJobParams struct {
	StudentID int32 `json:"student_id"`
	JobID     int32 `json:"job_id"`
}

func (q *Queries) ApplyForJob(ctx context.Context, arg ApplyForJobParams) error {
	_, err := q.db.Exec(ctx, applyForJob, arg.StudentID, arg.JobID)
	return err
}

const checkIfAlreadyAppliedForJob = `-- name: CheckIfAlreadyAppliedForJob :one
SELECT EXISTS(
    SELECT student_id
    FROM student_job_application_table
    WHERE job_id = $1
    AND student_id = $2
)
`

type CheckIfAlreadyAppliedForJobParams struct {
	JobID     int32 `json:"job_id"`
	StudentID int32 `json:"student_id"`
}

func (q *Queries) CheckIfAlreadyAppliedForJob(ctx context.Context, arg CheckIfAlreadyAppliedForJobParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkIfAlreadyAppliedForJob, arg.JobID, arg.StudentID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkIfAppliedForJobAlready = `-- name: CheckIfAppliedForJobAlready :one
SELECT EXISTS(
    SELECT job_id
    from student_job_application_table
    WHERE job_id = $1
    AND student_id = $2
)
`

type CheckIfAppliedForJobAlreadyParams struct {
	JobID     int32 `json:"job_id"`
	StudentID int32 `json:"student_id"`
}

func (q *Queries) CheckIfAppliedForJobAlready(ctx context.Context, arg CheckIfAppliedForJobAlreadyParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkIfAppliedForJobAlready, arg.JobID, arg.StudentID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getAlreadyAppliedJobIds = `-- name: GetAlreadyAppliedJobIds :many
SELECT job_id
FROM student_job_application_table
WHERE student_id = $1
`

func (q *Queries) GetAlreadyAppliedJobIds(ctx context.Context, studentID int32) ([]int32, error) {
	rows, err := q.db.Query(ctx, getAlreadyAppliedJobIds, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var job_id int32
		if err := rows.Scan(&job_id); err != nil {
			return nil, err
		}
		items = append(items, job_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAlreadyAppliedJobs = `-- name: GetAlreadyAppliedJobs :many
SELECT j.job_id, job_role, job_description, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches
FROM student_job_application_table as sj
JOIN job_table as j
ON sj.job_id = j.job_id
WHERE sj.student_id = $1
`

type GetAlreadyAppliedJobsRow struct {
	JobID            int32            `json:"job_id"`
	JobRole          string           `json:"job_role"`
	JobDescription   string           `json:"job_description"`
	JobType          string           `json:"job_type"`
	Ctc              float32          `json:"ctc"`
	SalaryTier       string           `json:"salary_tier"`
	ApplyByDate      pgtype.Timestamp `json:"apply_by_date"`
	CgpaCutoff       float32          `json:"cgpa_cutoff"`
	EligibleBatch    int32            `json:"eligible_batch"`
	EligibleBranches []string         `json:"eligible_branches"`
}

func (q *Queries) GetAlreadyAppliedJobs(ctx context.Context, studentID int32) ([]*GetAlreadyAppliedJobsRow, error) {
	rows, err := q.db.Query(ctx, getAlreadyAppliedJobs, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetAlreadyAppliedJobsRow
	for rows.Next() {
		var i GetAlreadyAppliedJobsRow
		if err := rows.Scan(
			&i.JobID,
			&i.JobRole,
			&i.JobDescription,
			&i.JobType,
			&i.Ctc,
			&i.SalaryTier,
			&i.ApplyByDate,
			&i.CgpaCutoff,
			&i.EligibleBatch,
			&i.EligibleBranches,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobOfferActByDate = `-- name: GetJobOfferActByDate :one
SELECT act_by_date 
FROM student_offer_table
WHERE student_id = $1 
AND job_id = $2
`

type GetJobOfferActByDateParams struct {
	StudentID int32 `json:"student_id"`
	JobID     int32 `json:"job_id"`
}

func (q *Queries) GetJobOfferActByDate(ctx context.Context, arg GetJobOfferActByDateParams) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, getJobOfferActByDate, arg.StudentID, arg.JobID)
	var act_by_date pgtype.Timestamp
	err := row.Scan(&act_by_date)
	return act_by_date, err
}

const getJobOffers = `-- name: GetJobOffers :many
SELECT job_table.job_id, company_name, job_role, job_type, ctc, salary_tier, action, action_date, act_by_date
FROM student_offer_table JOIN job_table 
ON student_offer_table.job_id = job_table.job_id
JOIN company_table
ON job_table.company_id = company_table.company_id
AND student_id = $1
`

type GetJobOffersRow struct {
	JobID       int32            `json:"job_id"`
	CompanyName string           `json:"company_name"`
	JobRole     string           `json:"job_role"`
	JobType     string           `json:"job_type"`
	Ctc         float32          `json:"ctc"`
	SalaryTier  string           `json:"salary_tier"`
	Action      string           `json:"action"`
	ActionDate  pgtype.Timestamp `json:"action_date"`
	ActByDate   pgtype.Timestamp `json:"act_by_date"`
}

func (q *Queries) GetJobOffers(ctx context.Context, studentID int32) ([]*GetJobOffersRow, error) {
	rows, err := q.db.Query(ctx, getJobOffers, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobOffersRow
	for rows.Next() {
		var i GetJobOffersRow
		if err := rows.Scan(
			&i.JobID,
			&i.CompanyName,
			&i.JobRole,
			&i.JobType,
			&i.Ctc,
			&i.SalaryTier,
			&i.Action,
			&i.ActionDate,
			&i.ActByDate,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getJobs = `-- name: GetJobs :many
SELECT job_id, job_role, job_description, ctc, salary_tier, apply_by_date, cgpa_cutoff, company_name, industry,
       (CASE WHEN
                 cgpa_cutoff <= (SELECT cgpa FROM student_table WHERE student_id = $1) THEN TRUE
             ELSE FALSE
           END) AS can_apply
FROM job_table JOIN company_table
                    ON job_table.company_id = company_table.company_id
WHERE (COALESCE(array_length($2::VARCHAR[], 1), 0) = 0 OR salary_tier = ANY($2))
  AND (COALESCE(array_length($3::VARCHAR[], 1), 0) = 0 OR salary_tier <> ANY($3))
  AND (COALESCE(array_length($4::VARCHAR[], 1), 0) = 0 OR job_role = ANY($4))
  AND (COALESCE(array_length($5::VARCHAR[], 1), 0) = 0 OR job_role <> ANY($5))
  AND (COALESCE(array_length($6::VARCHAR[], 1), 0) = 0 OR company_name = ANY($6))
  AND NOW() < apply_by_date
  AND (COALESCE(array_length($7::INT[], 1), 0) = 0 OR job_id <> ANY($7))
  AND ARRAY(SELECT branch FROM student_table WHERE student_id = $1) && eligible_branches
AND job_table.eligible_batch = (SELECT batch from student_table where student_id = $1)
`

type GetJobsParams struct {
	StudentID                 *int32   `json:"student_id"`
	SalaryTierFilter          []string `json:"salary_tier_filter"`
	DoNotShowSalaryTierFilter []string `json:"do_not_show_salary_tier_filter"`
	JobRoleFilter             []string `json:"job_role_filter"`
	DoNotShowJobTypeFilter    []string `json:"do_not_show_job_type_filter"`
	CompanyNameFilter         []string `json:"company_name_filter"`
	AlreadyAppliedJobID       []int32  `json:"already_applied_job_id"`
}

type GetJobsRow struct {
	JobID          int32            `json:"job_id"`
	JobRole        string           `json:"job_role"`
	JobDescription string           `json:"job_description"`
	Ctc            float32          `json:"ctc"`
	SalaryTier     string           `json:"salary_tier"`
	ApplyByDate    pgtype.Timestamp `json:"apply_by_date"`
	CgpaCutoff     float32          `json:"cgpa_cutoff"`
	CompanyName    string           `json:"company_name"`
	Industry       string           `json:"industry"`
	CanApply       bool             `json:"can_apply"`
}

func (q *Queries) GetJobs(ctx context.Context, arg GetJobsParams) ([]*GetJobsRow, error) {
	rows, err := q.db.Query(ctx, getJobs,
		arg.StudentID,
		arg.SalaryTierFilter,
		arg.DoNotShowSalaryTierFilter,
		arg.JobRoleFilter,
		arg.DoNotShowJobTypeFilter,
		arg.CompanyNameFilter,
		arg.AlreadyAppliedJobID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobsRow
	for rows.Next() {
		var i GetJobsRow
		if err := rows.Scan(
			&i.JobID,
			&i.JobRole,
			&i.JobDescription,
			&i.Ctc,
			&i.SalaryTier,
			&i.ApplyByDate,
			&i.CgpaCutoff,
			&i.CompanyName,
			&i.Industry,
			&i.CanApply,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOfferedJobInfo = `-- name: GetOfferedJobInfo :many
SELECT DISTINCT salary_tier, job_type
FROM job_table as j
JOIN student_offer_table as so
ON j.job_id = so.job_id
WHERE so.student_id = $1
`

type GetOfferedJobInfoRow struct {
	SalaryTier string `json:"salary_tier"`
	JobType    string `json:"job_type"`
}

func (q *Queries) GetOfferedJobInfo(ctx context.Context, studentID int32) ([]*GetOfferedJobInfoRow, error) {
	rows, err := q.db.Query(ctx, getOfferedJobInfo, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetOfferedJobInfoRow
	for rows.Next() {
		var i GetOfferedJobInfoRow
		if err := rows.Scan(&i.SalaryTier, &i.JobType); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStudentProfile = `-- name: GetStudentProfile :one
SELECT name, usn, branch, cgpa, batch, num_active_backlogs, email_id, counsellor_email_id
FROM student_table
WHERE usn = $1
`

type GetStudentProfileRow struct {
	Name              string  `json:"name"`
	Usn               string  `json:"usn"`
	Branch            string  `json:"branch"`
	Cgpa              float32 `json:"cgpa"`
	Batch             int32   `json:"batch"`
	NumActiveBacklogs int32   `json:"num_active_backlogs"`
	EmailID           string  `json:"email_id"`
	CounsellorEmailID string  `json:"counsellor_email_id"`
}

func (q *Queries) GetStudentProfile(ctx context.Context, usn string) (*GetStudentProfileRow, error) {
	row := q.db.QueryRow(ctx, getStudentProfile, usn)
	var i GetStudentProfileRow
	err := row.Scan(
		&i.Name,
		&i.Usn,
		&i.Branch,
		&i.Cgpa,
		&i.Batch,
		&i.NumActiveBacklogs,
		&i.EmailID,
		&i.CounsellorEmailID,
	)
	return &i, err
}

const performJobOfferAction = `-- name: PerformJobOfferAction :exec
UPDATE student_offer_table SET action = $3, action_date = NOW() 
WHERE student_id = $1
AND job_id = $2
`

type PerformJobOfferActionParams struct {
	StudentID int32  `json:"student_id"`
	JobID     int32  `json:"job_id"`
	Action    string `json:"action"`
}

func (q *Queries) PerformJobOfferAction(ctx context.Context, arg PerformJobOfferActionParams) error {
	_, err := q.db.Exec(ctx, performJobOfferAction, arg.StudentID, arg.JobID, arg.Action)
	return err
}
