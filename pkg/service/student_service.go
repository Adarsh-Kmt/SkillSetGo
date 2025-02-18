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
	AcceptJobOffer(studentId int, request entity.PerformJobOfferActionRequest) (httpError *helper.HTTPError)
	RejectJobOffer(studentId int, request entity.PerformJobOfferActionRequest) (httpError *helper.HTTPError)
	GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *helper.HTTPError)
	GetStudentProfile(studentId int) (profile *sqlc.GetStudentProfileRow, httpError *helper.HTTPError)
	GetAlreadyAppliedJobs(studentId int) (jobs []*sqlc.GetAlreadyAppliedJobsRow, httpError *helper.HTTPError)

	GetScheduledInterviews(studentId int) (interviews []*sqlc.GetInterviewsScheduledForStudentRow, httpError *helper.HTTPError)
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
		offeredJobList          []*sqlc.GetOfferedJobInfoRow
	)

	offeredJobList, err = db.Client.GetOfferedJobInfo(context.TODO(), int32(studentId))

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	hasInternship := false
	hasOpenDreamOffer := false
	hasDreamOffer := false

	for _, offeredJob := range offeredJobList {

		if offeredJob.SalaryTier == "Open Dream" {
			hasOpenDreamOffer = true
		}
		if offeredJob.SalaryTier == "Dream" {
			hasDreamOffer = true
		}
		if offeredJob.JobType == "Internship" {
			hasInternship = true
		}
	}

	if hasOpenDreamOffer && hasInternship {
		return nil, nil
	}

	studentIdParam := int32(studentId)

	if alreadyAppliedJobIdList, err = db.Client.GetAlreadyAppliedJobIds(context.TODO(), studentIdParam); err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	params := sqlc.GetJobsParams{
		StudentID:                 &studentIdParam,
		AlreadyAppliedJobID:       alreadyAppliedJobIdList,
		DoNotShowSalaryTierFilter: make([]string, 0),
		DoNotShowJobTypeFilter:    make([]string, 0),
	}

	if hasOpenDreamOffer {
		params.DoNotShowSalaryTierFilter = append(params.DoNotShowSalaryTierFilter, "Open Dream")
		params.DoNotShowSalaryTierFilter = append(params.DoNotShowSalaryTierFilter, "Dream")
	}
	if hasDreamOffer {
		params.DoNotShowSalaryTierFilter = append(params.DoNotShowSalaryTierFilter, "Dream")
	}
	if hasInternship {
		params.DoNotShowJobTypeFilter = append(params.DoNotShowJobTypeFilter, "Internship")
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

func (service *StudentServiceImpl) AcceptJobOffer(studentId int, request entity.PerformJobOfferActionRequest) (httpError *helper.HTTPError) {

	params1 := sqlc.GetJobOfferParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
	}
	jobOfferRow, err := db.Client.GetJobOffer(context.TODO(), params1)

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if time.Now().After(jobOfferRow.ActByDate.Time) {
		return &helper.HTTPError{StatusCode: 400, Error: "act_by_date expired"}
	}

	rejectOpenDreamOffers := false
	rejectDreamOffers := false
	rejectInternshipOffers := false

	if jobOfferRow.SalaryTier == "Dream" {
		rejectDreamOffers = true
	}
	if jobOfferRow.JobType == "Internship" {
		rejectInternshipOffers = true
	}
	if jobOfferRow.SalaryTier == "Open Dream" {
		rejectOpenDreamOffers = true
		rejectDreamOffers = true
	}

	params2 := sqlc.PerformJobOfferActionParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
		Action:    request.Action,
	}
	if err := db.Client.PerformJobOfferAction(context.TODO(), params2); err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	pendingJobOffers, err := db.Client.GetPendingOffers(context.TODO(), int32(studentId))

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	for _, pendingJobOffer := range pendingJobOffers {

		if pendingJobOffer.JobType == "Internship" && rejectInternshipOffers {

			rejectOfferParams := sqlc.RejectOfferParams{
				StudentID: int32(studentId),
				JobID:     pendingJobOffer.JobID,
			}
			err := db.Client.RejectOffer(context.TODO(), rejectOfferParams)
			if err != nil {
				return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
			}
		} else if pendingJobOffer.SalaryTier == "Dream" && rejectDreamOffers {

			rejectOfferParams := sqlc.RejectOfferParams{
				StudentID: int32(studentId),
				JobID:     pendingJobOffer.JobID,
			}
			err := db.Client.RejectOffer(context.TODO(), rejectOfferParams)
			if err != nil {
				return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
			}
		} else if pendingJobOffer.SalaryTier == "Open Dream" && rejectOpenDreamOffers {

			rejectOfferParams := sqlc.RejectOfferParams{
				StudentID: int32(studentId),
				JobID:     pendingJobOffer.JobID,
			}
			err := db.Client.RejectOffer(context.TODO(), rejectOfferParams)
			if err != nil {
				return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
			}
		}
	}
	return nil
}

func (service *StudentServiceImpl) RejectJobOffer(studentId int, request entity.PerformJobOfferActionRequest) (httpError *helper.HTTPError) {

	params1 := sqlc.GetJobOfferParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
	}
	jobOfferRow, err := db.Client.GetJobOffer(context.TODO(), params1)

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	if time.Now().After(jobOfferRow.ActByDate.Time) {
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

func (service *StudentServiceImpl) GetScheduledInterviews(studentId int) (interviews []*sqlc.GetInterviewsScheduledForStudentRow, httpError *helper.HTTPError) {

	rows, err := db.Client.GetInterviewsScheduledForStudent(context.TODO(), int32(studentId))

	if err != nil {
		return nil, &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return rows, nil
}
