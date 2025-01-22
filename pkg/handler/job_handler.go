package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
)

type JobHandler struct {
	js service.JobService
}

func NewJobHandler(js service.JobService) *JobHandler {
	return &JobHandler{js: js}
}
func (jh *JobHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/job", util.MakeHttpHandlerFunc(jh.GetJobs)).Methods("GET")
	mux.HandleFunc("/job", util.MakeHttpHandlerFunc(jh.CreateJob)).Methods("POST")
	return mux
}

func (jh *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	var (
		jobs []*sqlc.GetJobsRow
	)

	queryParams := r.URL.Query()
	salaryTierList := queryParams["salary-tier"]
	companyList := queryParams["company"]
	jobRoleList := queryParams["job-role"]

	for _, salaryTier := range salaryTierList {
		if salaryTier != "Dream" && salaryTier != "Open Dream" && salaryTier != "Mass Recruitment" {
			return &util.HTTPError{StatusCode: 400, Error: "invalid salary tier url query parameter"}
		}
	}

	//log.Println(salaryTierList)
	if jobs, httpError = jh.js.GetJobs(1, salaryTierList, jobRoleList, companyList); httpError != nil {
		return httpError
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
	return nil
}

// CreateJob TODO: implement functionality.
func (jh *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	request := &entity.CreateJobRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "bad request"}
	}
	log.Println(request)
	if httpError = entity.ValidateCreateJobRequest(*request); httpError != nil {
		return httpError
	}
	util.WriteJSON(w, http.StatusOK, map[string]any{"message": "job created successfully"})
	return nil
}
