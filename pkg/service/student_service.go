package service

import (
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type StudentService interface {
	RegisterStudent(request *entity.RegisterStudentRequest) (httpError *util.HTTPError)
}
