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
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/jackc/pgx/v5"
)

type StudentService interface {
	ApplyForJob(studentId int, jobId int) (httpError *util.HTTPError)
	GetJobOffers(studentId int) (offers []response.JobOfferResponse, httpError *util.HTTPError)
	PerformJobOfferAction(studentId int, request entity.PerformJobOfferActionRequest) (httpError *util.HTTPError)
}

type StudentServiceImpl struct {
}

func NewStudentServiceImpl() *StudentServiceImpl {
	return &StudentServiceImpl{}
}

func (ss *StudentServiceImpl) ApplyForJob(studentId int, jobId int) (httpError *util.HTTPError) {

	params := sqlc.RegisterForJobParams{
		StudentID: int32(studentId),
		JobID:     int32(jobId),
	}

	log.Println(params)

	if err := db.Client.RegisterForJob(context.TODO(), params); err != nil {
		log.Println(err.Error())
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	return nil
}

func (ss *StudentServiceImpl) GetJobOffers(studentId int) (offers []response.JobOfferResponse, httpError *util.HTTPError) {

	var err error

	offers = make([]response.JobOfferResponse, 0)

	offerRows, err := db.Client.GetJobOffers(context.Background(), int32(studentId))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return offers, nil
		}
		return nil, &util.HTTPError{StatusCode: 500, Error: "internal server error"}
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

func (ss *StudentServiceImpl) PerformJobOfferAction(studentId int, request entity.PerformJobOfferActionRequest) (httpError *util.HTTPError) {

	var (
		actByDate time.Time
	)
	params1 := sqlc.GetJobOfferActByDateParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
	}
	if actByDatePgxTimestamp, err := db.Client.GetJobOfferActByDate(context.TODO(), params1); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	} else {
		actByDate = actByDatePgxTimestamp.Time
	}

	if time.Now().After(actByDate) {
		return &util.HTTPError{StatusCode: 400, Error: "act_by_date expired"}
	}

	params2 := sqlc.PerformJobOfferActionParams{
		StudentID: int32(studentId),
		JobID:     int32(request.JobId),
		Action:    request.Action,
	}
	if err := db.Client.PerformJobOfferAction(context.TODO(), params2); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}
