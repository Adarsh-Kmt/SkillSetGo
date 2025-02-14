package service

import (
	"context"
	"log/slog"
	"time"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"github.com/jackc/pgx/v5/pgtype"
)

type CompanyService interface {
	CreateJob(companyId int, request entity.CreateJobRequest) (httpError *helper.HTTPError)
	OfferJob(request entity.OfferJobRequest) (httpError *helper.HTTPError)
	GetPublishedJobs(companyId int) (jobs []*sqlc.GetPublishedJobsRow, httpError *helper.HTTPError)
	GetJobApplicants(companyId int, jobId int) (profiles []*sqlc.GetJobApplicantsRow, httpError *helper.HTTPError)
}
type CompanyServiceImpl struct {
}

func NewCompanyServiceImpl() *CompanyServiceImpl {
	return &CompanyServiceImpl{}
}

func (service *CompanyServiceImpl) CreateJob(companyId int, request entity.CreateJobRequest) (httpError *helper.HTTPError) {

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
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return nil
}

func (service *CompanyServiceImpl) OfferJob(request entity.OfferJobRequest) (httpError *helper.HTTPError) {

	abd, _ := time.Parse("2006-01-02 15:04:05", request.ActByDate)

	params := sqlc.OfferJobParams{
		StudentID: int32(request.StudentId),
		JobID:     int32(request.JobId),
		ActByDate: pgtype.Timestamp{Time: abd, Valid: true},
	}

	err := db.Client.OfferJob(context.TODO(), params)

	if err != nil {
		slog.Error(err.Error())
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}

	}
	return nil
}

func (service *CompanyServiceImpl) GetPublishedJobs(companyId int) (jobs []*sqlc.GetPublishedJobsRow, httpError *helper.HTTPError) {

	rows, err := db.Client.GetPublishedJobs(context.TODO(), int32(companyId))

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return rows, nil
}

func (service *CompanyServiceImpl) GetJobApplicants(companyId int, jobId int) (profiles []*sqlc.GetJobApplicantsRow, httpError *helper.HTTPError) {

	checkParams := sqlc.CheckIfCompanyCreatedJobParams{CompanyID: int32(companyId), JobID: int32(jobId)}

	wasCreated, err := db.Client.CheckIfCompanyCreatedJob(context.TODO(), checkParams)

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if !wasCreated {
		return nil, &helper.HTTPError{StatusCode: 403, Error: "company did not publish job"}
	}

	rows, err := db.Client.GetJobApplicants(context.TODO(), int32(jobId))

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return rows, nil
}
