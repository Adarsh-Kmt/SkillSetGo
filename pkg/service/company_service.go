package service

import (
	"context"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type CompanyService interface {
	RegisterCompany(CompanyName string, PocName string, PocPhno string, Industry string) (httpError *util.HTTPError)
}

type CompanyServiceImpl struct {
}

func NewCompanyServiceImpl() *CompanyServiceImpl {
	return &CompanyServiceImpl{}
}

func (cs *CompanyServiceImpl) RegisterCompany(CompanyName string, PocName string, PocPhno string, Industry string) (httpError *util.HTTPError) {
	var (
		err error
	)
	params := sqlc.InsertCompanyParams{
		CompanyName: CompanyName,
		PocName:     PocName,
		PocPhno:     PocPhno,
		Industry:    Industry,
	}

	if err = db.Client.InsertCompany(context.TODO(), params); err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return nil
}
