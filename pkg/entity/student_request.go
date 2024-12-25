package entity

import (
	"strings"

	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type RegisterStudentRequest struct {
	Name   string  `json:"name"`
	usn    string  `json:"usn"`
	Branch string  `json:"branch"`
	Cgpa   float32 `json:"cgpa"`
	Year   int     `json:"year"`
	Email  string  `json:"email"`
	Phone  string  `json:"phone"`
}

func ValidateRegisterStudentRequest(request RegisterStudentRequest) (httpError *util.HTTPError) {
	if request.Name == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Name cannot be empty"}
	}
	if !strings.HasPrefix(request.usn, "1RV") {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid USN"}
	}
	if request.Branch == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Branch cannot be empty"}
	}
	if request.Cgpa <= 0 || request.Cgpa > 10 {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid CGPA"}
	}
	if request.Year == 0 || request.Year > 4 {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Year"}
	}
	if strings.HasSuffix(request.Email, "@rvce.edu.in") {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Email ID"}
	}
	if request.Phone == "" || len(request.Phone) != 10 {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Phone Number"}
	}
	return nil
}
