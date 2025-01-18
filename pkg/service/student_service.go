package service

import (
	"context"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type StudentService interface {
	RegisterStudent(request entity.RegisterStudentRequest) (httpError *util.HTTPError)
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
		Usn:               request.Usn,
		Sname:             request.Name,
		Branch:            request.Branch,
		Cgpa:              request.Cgpa,
		EmailID:           request.Email,
		PhoneNumber:       request.Phone,
		CounsellorEmailID: request.CounsellorEmailID,
		NumActiveBacklogs: int32(request.NumberOfBacklogs),
	}

	if err = db.Client.InsertUser(context.TODO(), params); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}
