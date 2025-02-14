package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/adarsh-kmt/skillsetgo/pkg/response"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"github.com/jackc/pgx/v5"
)

type StudentService interface {
	ApplyForJob(studentId int, jobId int) (httpError *helper.HTTPError)
	GetJobOffers(studentId int) (offers []response.JobOfferResponse, httpError *helper.HTTPError)
	PerformJobOfferAction(studentId int, request entity.PerformJobOfferActionRequest) (httpError *helper.HTTPError)
	GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *helper.HTTPError)
	GetStudentProfile(studentId int) (profile *sqlc.GetStudentProfileRow, httpError *helper.HTTPError)
	GetAlreadyAppliedJobs(studentId int) (jobs []*sqlc.GetAlreadyAppliedJobsRow, httpError *helper.HTTPError)
}

type StudentServiceImpl struct {
}

func NewStudentServiceImpl() *StudentServiceImpl {
	return &StudentServiceImpl{}
}

func (service *StudentServiceImpl) GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *helper.HTTPError) {

	var (
		err                     error
		alreadyAppliedJobIdList []int32
	)
	studentIdParam := int32(studentId)

	if alreadyAppliedJobIdList, err = db.Client.GetAlreadyAppliedJobIds(context.TODO(), studentIdParam); err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	params := sqlc.GetJobsParams{
		Column1:             &studentIdParam,
		Column2:             salaryTierFilter,
		Column3:             jobRoleFilter,
		Column4:             companyFilter,
		AlreadyAppliedJobID: alreadyAppliedJobIdList,
	}

	if jobs, err = db.Client.GetJobs(context.Background(), params); err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return jobs, nil
}

func (service *StudentServiceImpl) ApplyForJob(studentId int, jobId int) (httpError *helper.HTTPError) {

	checkParams := sqlc.CheckIfAppliedForJobAlreadyParams{
		StudentID: int32(studentId),
		JobID:     int32(jobId),
	}

	alreadyApplied, err := db.Client.CheckIfAppliedForJobAlready(context.TODO(), checkParams)

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if alreadyApplied {
		return &helper.HTTPError{StatusCode: 404, Error: "already applied for job"}
	}
	params := sqlc.ApplyForJobParams{
		StudentID: int32(studentId),
		JobID:     int32(jobId),
	}

	if err := db.Client.ApplyForJob(context.TODO(), params); err != nil {
		log.Println(err.Error())
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return nil
}

func (service *StudentServiceImpl) GetJobOffers(studentId int) (offers []response.JobOfferResponse, httpError *helper.HTTPError) {

	var err error

	offers = make([]response.JobOfferResponse, 0)

	offerRows, err := db.Client.GetJobOffers(context.Background(), int32(studentId))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return offers, nil
		}
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	for _, offerRow := range offerRows {

		offer := response.JobOfferResponse{
			JobId:       offerRow.JobID,
			CompanyName: offerRow.CompanyName,
			JobRole:     offerRow.JobRole,
			JobType:     offerRow.JobType,
			CTC:         offerRow.Ctc,
			SalaryTier:  offerRow.SalaryTier,
			Action:      offerRow.Action,
			ActionDate:  offerRow.ActionDate.Time.String(),
			ActByDate:   offerRow.ActByDate.Time.String(),
		}

		offers = append(offers, offer)

	}
	return offers, nil
}

func (service *StudentServiceImpl) PerformJobOfferAction(studentId int, request entity.PerformJobOfferActionRequest) (httpError *helper.HTTPError) {

	var (
		actByDate time.Time
	)
	params1 := sqlc.GetJobOfferActByDateParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
	}
	if actByDatePgxTimestamp, err := db.Client.GetJobOfferActByDate(context.TODO(), params1); err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	} else {
		actByDate = actByDatePgxTimestamp.Time
	}

	if time.Now().After(actByDate) {
		return &helper.HTTPError{StatusCode: 400, Error: "act_by_date expired"}
	}

	params2 := sqlc.PerformJobOfferActionParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
		Action:    request.Action,
	}
	if err := db.Client.PerformJobOfferAction(context.TODO(), params2); err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}

func (service *StudentServiceImpl) GetStudentProfile(studentId int) (profile *sqlc.GetStudentProfileRow, httpError *helper.HTTPError) {

	row, err := db.Client.GetStudentProfile(context.TODO(), int32(studentId))
	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return row, nil
}

func (service *StudentServiceImpl) GetAlreadyAppliedJobs(studentId int) (jobs []*sqlc.GetAlreadyAppliedJobsRow, httpError *helper.HTTPError) {

	rows, err := db.Client.GetAlreadyAppliedJobs(context.TODO(), int32(studentId))
	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return rows, nil
}
