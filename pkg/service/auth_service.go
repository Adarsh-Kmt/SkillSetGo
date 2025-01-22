package service

import (
	"context"
	"errors"
	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/jackc/pgx/v5"
)

type AuthService interface {
	LoginStudent(request *entity.LoginStudentRequest) (accessToken string, httpError *util.HTTPError)
	LoginCompany(request *entity.LoginCompanyRequest) (accessToken string, httpError *util.HTTPError)
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
