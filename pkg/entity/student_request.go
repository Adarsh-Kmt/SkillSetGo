package entity

import "github.com/adarsh-kmt/skillsetgo/pkg/util"

type LoginStudentRequest struct {
	USN      string `json:"usn"`
	Password string `json:"password"`
}

type RegisterStudentRequest struct {
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
