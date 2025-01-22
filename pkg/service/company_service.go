package service

import (
	"context"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type CompanyService interface {
	RegisterCompany(request entity.RegisterCompanyRequest) (httpError *util.HTTPError)
}

type CompanyServiceImpl struct {
}

func NewCompanyServiceImpl() *CompanyServiceImpl {
	return &CompanyServiceImpl{}
}

func (cs *CompanyServiceImpl) RegisterCompany(request entity.RegisterCompanyRequest) (httpError *util.HTTPError) {
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
