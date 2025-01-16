package entity

import (
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type RegisterCompanyRequest struct {
	CompanyName string `json:"company_name"`
	PocName     string `json:"poc_name"`
	PocPhno     string `json:"poc_phno"`
	Industry    string `json:"industry"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func ValidateRegisterCompanyRequest(request RegisterCompanyRequest) (httpError *util.HTTPError) {
	if request.CompanyName == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Company Name cannot be empty"}
	}
	if request.PocName == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Point of Contact Name cannot be empty"}
	}
	if request.PocPhno == "" || len(request.PocPhno) != 10 {
		return &util.HTTPError{StatusCode: 400, Error: "Invalid Phone Number"}
	}
	if request.Industry == "" {
		return &util.HTTPError{StatusCode: 400, Error: "Industry cannot be empty"}
	}
	return nil
}
