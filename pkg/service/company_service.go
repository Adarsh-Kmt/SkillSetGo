package service

import (
	"context"
	"fmt"
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
	GetOfferStatus(companyId int, jobId int) (offerStatus []*sqlc.GetOfferStatusRow, httpError *helper.HTTPError)

	ScheduleInterview(companyId int, request entity.ScheduleInterviewRequest) (httpError *helper.HTTPError)
	GetScheduledInterviews(companyId int, jobId int) (interviews []*sqlc.GetInterviewsScheduledByCompanyRow, httpError *helper.HTTPError)

	GetPlacementStats() (stats []*sqlc.GetPlacementStatsRow, httpError *helper.HTTPError)
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
		JobDescription:   request.JobDescription,
	}
	if err := db.Client.CreateJob(context.TODO(), params); err != nil {
		slog.Error(err.Error())
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return nil
}

func (service *CompanyServiceImpl) OfferJob(request entity.OfferJobRequest) (httpError *helper.HTTPError) {

	abd, _ := time.Parse("2006-01-02 15:04:05", request.ActByDate)

	checkParams := sqlc.CheckIfOfferedAlreadyParams{
		StudentID: int32(request.StudentId),
		JobID:     int32(request.JobId),
	}

	exists, err := db.Client.CheckIfOfferedAlready(context.TODO(), checkParams)

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if exists {
		return &helper.HTTPError{StatusCode: 404, Error: "student already received offer"}
	}
	params := sqlc.OfferJobParams{
		StudentID: int32(request.StudentId),
		JobID:     int32(request.JobId),
		ActByDate: pgtype.Timestamp{Time: abd, Valid: true},
	}

	err = db.Client.OfferJob(context.TODO(), params)

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

func (service *CompanyServiceImpl) GetOfferStatus(companyId int, jobId int) (offerStatus []*sqlc.GetOfferStatusRow, httpError *helper.HTTPError) {

	checkParams := sqlc.CheckIfCompanyCreatedJobParams{CompanyID: int32(companyId), JobID: int32(jobId)}

	wasCreated, err := db.Client.CheckIfCompanyCreatedJob(context.TODO(), checkParams)

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if !wasCreated {
		return nil, &helper.HTTPError{StatusCode: 403, Error: "company did not publish job"}
	}

	rows, err := db.Client.GetOfferStatus(context.TODO(), int32(jobId))

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return rows, nil
}

func (service *CompanyServiceImpl) ScheduleInterview(companyId int, request entity.ScheduleInterviewRequest) (httpError *helper.HTTPError) {

	checkParams := sqlc.CheckIfCompanyCreatedJobParams{CompanyID: int32(companyId), JobID: int32(request.JobId)}

	wasCreated, err := db.Client.CheckIfCompanyCreatedJob(context.TODO(), checkParams)

	if err != nil {
		slog.Error(fmt.Sprintf("error while checking if company created job %s", err.Error()))
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if !wasCreated {
		return &helper.HTTPError{StatusCode: 403, Error: "company did not publish job"}
	}

	interviewDate, _ := time.Parse("2006-01-02 15:04:05", request.InterviewDate)

	existsCheckParams := sqlc.CheckIfInterviewScheduledAlreadyParams{
		StudentID:     int32(request.StudentId),
		JobID:         int32(request.JobId),
		InterviewDate: pgtype.Timestamp{Time: interviewDate, Valid: true},
		Venue:         request.Venue,
	}

	exists, err := db.Client.CheckIfInterviewScheduledAlready(context.TODO(), existsCheckParams)

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if exists {
		return &helper.HTTPError{StatusCode: 404, Error: "interview scheduled already"}
	}

	venueCheckParams := sqlc.CheckIfVenueBeingUsedAtParticularTimeParams{
		Venue:                        request.Venue,
		DateOfInterviewToBeScheduled: pgtype.Timestamp{Time: interviewDate, Valid: true},
	}
	exists, err = db.Client.CheckIfVenueBeingUsedAtParticularTime(context.TODO(), venueCheckParams)

	if err != nil {
		slog.Error(fmt.Sprintf("check venue error: %s", err.Error()))
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	if exists {
		return &helper.HTTPError{StatusCode: 404, Error: "venue already in use, schedule for another time"}
	}

	params := sqlc.ScheduleInterviewParams{
		StudentID:     int32(request.StudentId),
		JobID:         int32(request.JobId),
		Venue:         request.Venue,
		InterviewDate: pgtype.Timestamp{Time: interviewDate, Valid: true},
	}

	if err := db.Client.ScheduleInterview(context.TODO(), params); err != nil {
		slog.Error(fmt.Sprintf("error while scheduling interview: %s", err.Error()))
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return nil
}
func (service *CompanyServiceImpl) GetScheduledInterviews(companyId int, jobId int) (interviews []*sqlc.GetInterviewsScheduledByCompanyRow, httpError *helper.HTTPError) {

	checkParams := sqlc.CheckIfCompanyCreatedJobParams{CompanyID: int32(companyId), JobID: int32(jobId)}

	wasCreated, err := db.Client.CheckIfCompanyCreatedJob(context.TODO(), checkParams)

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if !wasCreated {
		return nil, &helper.HTTPError{StatusCode: 403, Error: "company did not publish job"}
	}

	rows, err := db.Client.GetInterviewsScheduledByCompany(context.TODO(), int32(jobId))
	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return rows, nil
}

func (service *CompanyServiceImpl) GetPlacementStats() (stats []*sqlc.GetPlacementStatsRow, httpError *helper.HTTPError) {

	rows, err := db.Client.GetPlacementStats(context.TODO())

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return rows, nil
}
