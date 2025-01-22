package entity

import (
	"strings"

	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type LoginStudentRequest struct {
	USN      string `json:"usn"`
	Password string `json:"password"`
}

type RegisterStudentRequest struct {
	Name              string  `json:"name"`
	Usn               string  `json:"usn"`
	Password          string  `json:"password"`
	Branch            string  `json:"branch"`
	Cgpa              float32 `json:"cgpa"`
	Email             string  `json:"email"`
	Batch             int     `json:"batch"`
	CounsellorEmailID string  `json:"counsellor_email_id"`
	NumberOfBacklogs  int     `json:"num_of_backlogs"`
}

type LoginCompanyRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ValidateLoginStudentRequest(request LoginStudentRequest) *util.HTTPError {
	errorMap := make(map[string]any)

	if request.USN == "" {
		errorMap["usn"] = "usn field cannot be empty"
	}
	if request.Password == "" {
		errorMap["password"] = "password field cannot be empty"
	}

	if len(errorMap) > 0 {
		return &util.HTTPError{Error: errorMap, StatusCode: 400}
	}
	return nil
}

func ValidateLoginCompanyRequest(request LoginCompanyRequest) *util.HTTPError {
	errorMap := make(map[string]any)

	if request.Username == "" {
		errorMap["usn"] = "usn field cannot be empty"
	}
	if request.Password == "" {
		errorMap["password"] = "password field cannot be empty"
	}

	if len(errorMap) > 0 {
		return &util.HTTPError{Error: errorMap, StatusCode: 400}
	}
	return nil
}

func ValidateRegisterStudentRequest(request RegisterStudentRequest) (httpError *util.HTTPError) {

	branchemail := map[string]string{
		"ISE":  ".is",
		"CSE":  ".cs",
		"CD":   ".cd",
		"CY":   ".cy",
		"AIML": ".ai",
		"ECE":  ".ec",
		"EEE":  ".ee",
		"ETE":  ".et",
		"EIE":  ".ei",
		"ME":   ".me",
		"CV":   ".cv",
		"BT":   ".bt",
		"CH":   ".ch",
		"IEM":  ".iem",
		"ASE":  ".ae",
	}
	if request.Name == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Name cannot be empty"}
	}
	if request.Password == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Password cannot be empty"}
	}
	if !strings.HasPrefix(request.Usn, "1RV") {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid USN"}
	}
	substr, exists := branchemail[request.Branch]
	if request.Branch == "" || !exists {
		return &util.HTTPError{StatusCode: 400, Error: "Branch is empty or invalid"}
	}
	if request.Cgpa <= 0 || request.Cgpa > 10 {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid CGPA"}
	}
	email := strings.TrimSpace(request.Email)
	if !strings.HasSuffix(strings.ToLower(email), "@rvce.edu.in") {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Email ID"}
	}
	if !strings.Contains(request.Email, substr) {
		return &util.HTTPError{StatusCode: 400, Error: "Email ID does not match branch"}
	}
	if request.Batch < 2026 {
		return &util.HTTPError{StatusCode: 400, Error: "invalid batch"}
	}
	if request.CounsellorEmailID == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Counsellor Email ID cannot be empty"}
	}
	return nil
}
