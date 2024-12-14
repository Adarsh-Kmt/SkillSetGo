package service

import (
	"context"
	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
)

type JobService interface {
	GetJobs(studentId int, salaryTierFilter []string) (jobs []*sqlc.GetJobsRow, httpError *util.HTTPError)
}
type JobServiceImpl struct {
}

func NewJobServiceImpl() *JobServiceImpl {
	return &JobServiceImpl{}
}
func (js *JobServiceImpl) GetJobs(studentId int, salaryTierFilter []string) (jobs []*sqlc.GetJobsRow, httpError *util.HTTPError) {

	var (
		err error
	)
	params := sqlc.GetJobsParams{
		StudentID: int32(studentId),
		Column2:   salaryTierFilter,
	}

	if jobs, err = db.Client.GetJobs(context.Background(), params); err != nil {
		return nil, &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}
	return jobs, nil
}
