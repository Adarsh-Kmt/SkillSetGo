package service

import (
	"context"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type StudentService interface {
	RegisterStudent(Name []string, Branch []string, CGPA float64, EmailID []string, USN []string, ActiveBacklogs bool, CounsellorName []string) (httpError *util.HTTPError)
}

type StudentServiceImpl struct {
}

func NewStudentServiceImpl() *StudentServiceImpl {
	return &StudentServiceImpl{}
}

func (ss *StudentServiceImpl) RegisterStudent(request entity.RegisterStudentRequest) (httpError *util.HTTPError) {
	var (
		err error
	)
	params := sqlc.InsertUserParams{
		Usn:     request.Usn,
		Name:    request.Name,
		Branch:  request.Branch,
		Cgpa:    request.Cgpa,
		EmailID: request.Email,
	}

	if err = db.Client.InsertUser(context.TODO(), params); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

}
