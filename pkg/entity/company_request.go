package entity

import (
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"time"
)

type RegisterCompanyRequest struct {
	CompanyName string `json:"company_name"`
	PocName     string `json:"poc_name"`
	PocPhno     string `json:"poc_phno"`
	Industry    string `json:"industry"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type ScheduleInterviewRequest struct {
	JobId          int    `json:"job_id"`
	StudentId      int    `json:"student_id"`
	Venue          string `json:"venue"`
	InterviewDate  string `json:"interview_date"`
	InterviewRound string `json:"interview_round"`
}

type UpdateInterviewResultRequest struct {
	JobId          int    `json:"job_id"`
	StudentId      int    `json:"student_id"`
	InterviewRound int    `json:"interview_round"`
	Result         string `json:"result"`
}

func ValidateRegisterCompanyRequest(request RegisterCompanyRequest) (httpError *helper.HTTPError) {
	if request.CompanyName == "" {
		return &helper.HTTPError{StatusCode: 400, Error: "Company Name cannot be empty"}
	}
	if request.PocName == "" {
		return &helper.HTTPError{StatusCode: 400, Error: "Point of Contact Name cannot be empty"}
	}
	if request.PocPhno == "" || len(request.PocPhno) != 10 {
		return &helper.HTTPError{StatusCode: 400, Error: "Invalid Phone Number"}
	}
	if request.Industry == "" {
		return &helper.HTTPError{StatusCode: 400, Error: "Industry cannot be empty"}
	}
	return nil
}

func ValidateScheduleInterviewRequest(request ScheduleInterviewRequest) (httpError *helper.HTTPError) {

	if request.JobId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "Job Id cannot be empty"}
	}
	if request.StudentId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "Student Id cannot be empty"}
	}
	if request.Venue == "" {
		return &helper.HTTPError{StatusCode: 400, Error: "Venue cannot be empty"}
	}
	interviewDate, err := time.Parse("2006-01-02 15:04:05", request.InterviewDate)

	if err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "Invalid Interview Date format"}
	}

	if time.Now().After(interviewDate) {
		return &helper.HTTPError{StatusCode: 400, Error: "Interview Date expired"}
	}

	return nil
}

func ValidateUpdateInterviewResultRequest(request UpdateInterviewResultRequest) (httpError *helper.HTTPError) {

	if request.JobId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "Job Id cannot be empty"}
	}
	if request.StudentId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "Student Id cannot be empty"}
	}
	if request.InterviewRound == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "Interview Round cannot be 0"}
	}
	if request.Result != "PASSED" && request.Result != "FAILED" {
		return &helper.HTTPError{StatusCode: 400, Error: "Invalid Result"}
	}

	return nil
}
