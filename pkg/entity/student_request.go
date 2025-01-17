package entity

import (
	"strings"

	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type RegisterStudentRequest struct {
	Name                string  `json:"name"`
	Usn                 string  `json:"usn"`
	Branch              string  `json:"branch"`
	Cgpa                float32 `json:"cgpa"`
	Email               string  `json:"email"`
	Phone               string  `json:"phone"`
	counsellor_email_id string  `json:"counsellor_email_id"`
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

	if strings.HasSuffix(request.Email, "@rvce.edu.in") {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Email ID"}
	}
	if !strings.Contains(request.Email, substr) {
		return &util.HTTPError{StatusCode: 400, Error: "Email ID does not match branch"}
	}
	if request.Phone == "" || len(request.Phone) != 10 {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Phone Number"}
	}
	return nil
}
