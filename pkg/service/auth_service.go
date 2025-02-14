package service

import (
	"context"
	"errors"
	"log"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"github.com/jackc/pgx/v5"
)

type AuthService interface {
	LoginStudent(request *entity.LoginStudentRequest) (accessToken string, httpError *helper.HTTPError)
	LoginCompany(request *entity.LoginCompanyRequest) (accessToken string, httpError *helper.HTTPError)
	RegisterStudent(request *entity.RegisterStudentRequest) (httpError *helper.HTTPError)
	RegisterCompany(request *entity.RegisterCompanyRequest) (httpError *helper.HTTPError)
}

type AuthServiceImpl struct {
}

func NewAuthServiceImpl() *AuthServiceImpl {
	return &AuthServiceImpl{}
}
func (as *AuthServiceImpl) LoginStudent(request *entity.LoginStudentRequest) (accessToken string, httpError *helper.HTTPError) {
	params := sqlc.AuthenticateStudentParams{
		Usn:      request.USN,
		Password: request.Password,
	}
	studentId, err := db.Client.AuthenticateStudent(context.TODO(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &helper.HTTPError{StatusCode: 401, Error: "invalid username or password"}
		} else {
			return "", &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
		}
	}

	accessToken, httpError = helper.IssueToken(studentId, []string{"student"})

	return accessToken, httpError
}

func (as *AuthServiceImpl) LoginCompany(request *entity.LoginCompanyRequest) (accessToken string, httpError *helper.HTTPError) {
	params := sqlc.AuthenticateCompanyParams{
		Username: request.Username,
		Password: request.Password,
	}
	companyId, err := db.Client.AuthenticateCompany(context.TODO(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &helper.HTTPError{StatusCode: 401, Error: "invalid username or password"}
		} else {
			return "", &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
		}
	}

	accessToken, httpError = helper.IssueToken(companyId, []string{"company admin"})

	return accessToken, httpError
}

func (as *AuthServiceImpl) RegisterStudent(request *entity.RegisterStudentRequest) (httpError *helper.HTTPError) {

	var (
		err error
	)

	if exists, err := db.Client.CheckIfStudentExists(context.TODO(), request.Usn); err != nil {

		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	} else if exists {
		return &helper.HTTPError{StatusCode: 404, Error: "student already registered"}
	}
	params := sqlc.InsertUserParams{
		Usn:               request.Usn,
		Name:              request.Name,
		Password:          request.Password,
		Branch:            request.Branch,
		Batch:             int32(request.Batch),
		Cgpa:              request.Cgpa,
		EmailID:           request.Email,
		CounsellorEmailID: request.CounsellorEmailID,
		NumActiveBacklogs: int32(request.NumberOfBacklogs),
	}

	if err = db.Client.InsertUser(context.TODO(), params); err != nil {
		log.Println(err.Error())
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}

func (as *AuthServiceImpl) RegisterCompany(request *entity.RegisterCompanyRequest) (httpError *helper.HTTPError) {

	var (
		err error
	)

	if exists, err := db.Client.CheckIfCompanyExists(context.TODO(), request.CompanyName); err != nil {

		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	} else if exists {
		return &helper.HTTPError{StatusCode: 404, Error: "company already registered"}
	}
	params := sqlc.CreateCompanyParams{
		CompanyName: request.CompanyName,
		PocName:     request.PocName,
		PocPhno:     request.PocPhno,
		Industry:    request.Industry,
		Username:    request.Username,
		Password:    request.Password,
	}

	if err = db.Client.CreateCompany(context.TODO(), params); err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}
