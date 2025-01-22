package service

import (
	"context"
	"errors"
	"log"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/jackc/pgx/v5"
)

type AuthService interface {
	LoginStudent(request *entity.LoginStudentRequest) (accessToken string, httpError *util.HTTPError)
	LoginCompany(request *entity.LoginCompanyRequest) (accessToken string, httpError *util.HTTPError)
	RegisterStudent(request *entity.RegisterStudentRequest) (httpError *util.HTTPError)
	RegisterCompany(request *entity.RegisterCompanyRequest) (httpError *util.HTTPError)
}

type AuthServiceImpl struct {
}

func NewAuthServiceImpl() *AuthServiceImpl {
	return &AuthServiceImpl{}
}
func (as *AuthServiceImpl) LoginStudent(request *entity.LoginStudentRequest) (accessToken string, httpError *util.HTTPError) {
	params := sqlc.AuthenticateStudentParams{
		Usn:      request.USN,
		Password: request.Password,
	}
	studentId, err := db.Client.AuthenticateStudent(context.TODO(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &util.HTTPError{StatusCode: 401, Error: "invalid username or password"}
		} else {
			return "", &util.HTTPError{StatusCode: 500, Error: "internal server error"}
		}
	}

	accessToken, httpError = util.IssueToken(studentId)

	return accessToken, httpError
}

func (as *AuthServiceImpl) LoginCompany(request *entity.LoginCompanyRequest) (accessToken string, httpError *util.HTTPError) {
	params := sqlc.AuthenticateCompanyParams{
		Username: request.Username,
		Password: request.Password,
	}
	companyId, err := db.Client.AuthenticateCompany(context.TODO(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &util.HTTPError{StatusCode: 401, Error: "invalid username or password"}
		} else {
			return "", &util.HTTPError{StatusCode: 500, Error: "internal server error"}
		}
	}

	accessToken, httpError = util.IssueToken(companyId)

	return accessToken, httpError
}

func (as *AuthServiceImpl) RegisterStudent(request *entity.RegisterStudentRequest) (httpError *util.HTTPError) {

	var (
		err error
	)
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
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}

func (as *AuthServiceImpl) RegisterCompany(request *entity.RegisterCompanyRequest) (httpError *util.HTTPError) {

	var (
		err error
	)
	params := sqlc.CreateCompanyParams{
		CompanyName: request.CompanyName,
		PocName:     request.PocName,
		PocPhno:     request.PocPhno,
		Industry:    request.Industry,
		Username:    request.Username,
		Password:    request.Password,
	}

	if err = db.Client.CreateCompany(context.TODO(), params); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}
