package service

import (
	"context"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type JobService interface {
	GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *util.HTTPError)
}
type JobServiceImpl struct {
}

func NewJobServiceImpl() *JobServiceImpl {
	return &JobServiceImpl{}
}
func (js *JobServiceImpl) GetJobs(studentId int, salaryTierFilter []string, jobRoleFilter []string, companyFilter []string) (jobs []*sqlc.GetJobsRow, httpError *util.HTTPError) {

	var (
		err error
	)
	studentIdParam := int32(studentId)
	params := sqlc.GetJobsParams{
		Column1: &studentIdParam,
		Column2: salaryTierFilter,
		Column3: jobRoleFilter,
		Column4: companyFilter,
	}

	if jobs, err = db.Client.GetJobs(context.Background(), params); err != nil {
		return nil, &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return jobs, nil
}
