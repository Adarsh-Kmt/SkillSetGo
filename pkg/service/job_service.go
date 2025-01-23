package service

import (
	"context"
	"log/slog"
	"time"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobService interface {
	GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *util.HTTPError)
	CreateJob(companyId int, request entity.CreateJobRequest) (httpError *util.HTTPError)
	OfferJob(request entity.OfferJobRequest) (httpError *util.HTTPError)
}
type JobServiceImpl struct {
}

func NewJobServiceImpl() *JobServiceImpl {
	return &JobServiceImpl{}
}
func (js *JobServiceImpl) GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *util.HTTPError) {

	var (
		err                     error
		alreadyAppliedJobIdList []int32
	)
	studentIdParam := int32(studentId)

	if alreadyAppliedJobIdList, err = db.Client.GetAlreadyAppliedJobs(context.TODO(), studentIdParam); err != nil {
		return nil, &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	params := sqlc.GetJobsParams{
		Column1:             &studentIdParam,
		Column2:             salaryTierFilter,
		Column3:             jobRoleFilter,
		Column4:             companyFilter,
		AlreadyAppliedJobID: alreadyAppliedJobIdList,
	}

	if jobs, err = db.Client.GetJobs(context.Background(), params); err != nil {
		return nil, &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return jobs, nil
}

func (js *JobServiceImpl) CreateJob(companyId int, request entity.CreateJobRequest) (httpError *util.HTTPError) {

	var abd time.Time

	abd, _ = time.Parse("2006-01-02 15:04:05", request.ApplyByDate)

	params := sqlc.CreateJobParams{
		CompanyID:        int32(companyId),
		JobRole:          request.JobRole,
		JobType:          request.JobType,
		Ctc:              request.Ctc,
		SalaryTier:       request.SalaryTier,
		ApplyByDate:      pgtype.Timestamp{Time: abd, Valid: true},
		CgpaCutoff:       request.CgpaCutoff,
		EligibleBatch:    int32(request.EligibleBatch),
		EligibleBranches: request.EligibleBranches,
	}
	if err := db.Client.CreateJob(context.TODO(), params); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return nil
}

func (js *JobServiceImpl) OfferJob(request entity.OfferJobRequest) (httpError *util.HTTPError) {

	abd, _ := time.Parse("2006-01-02 15:04:05", request.ActByDate)

	params := sqlc.OfferJobParams{
		StudentID: int32(request.StudentId),
		JobID:     int32(request.JobId),
		ActByDate: pgtype.Timestamp{Time: abd, Valid: true},
	}

	err := db.Client.OfferJob(context.TODO(), params)

	if err != nil {
		slog.Error(err.Error())
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}

	}
	return nil
}
