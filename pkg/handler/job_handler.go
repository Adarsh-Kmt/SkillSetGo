package handler

import (
	"encoding/json"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
	"net/http"
)

type JobHandler struct {
	jobService service.JobService
}

func NewJobHandler(js service.JobService) *JobHandler {
	return &JobHandler{jobService: js}
}
func (jh *JobHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/job", util.MakeHttpHandlerFunc(jh.GetJobs)).Methods("GET")
	mux.HandleFunc("/job", util.MakeHttpHandlerFunc(jh.CreateJob)).Methods("POST")
	mux.HandleFunc("/job/offer", util.MakeAuthenticatedHandler(util.MakeHttpHandlerFunc(jh.OfferJob))).Methods("POST")
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
	if jobs, httpError = jh.jobService.GetJobs(1, salaryTierList, jobRoleList, companyList); httpError != nil {
		return httpError
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
	return nil
}

// CreateJob TODO: implement functionality.
func (jh *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	companyId, err := util.ValidateAccessToken(r.Header.Get("Auth"))

	if err != nil {
		return &util.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	request := entity.CreateJobRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "bad request"}
	}
	//log.Println(request)
	if httpError = entity.ValidateCreateJobRequest(request); httpError != nil {
		return httpError
	}

	if httpError = jh.jobService.CreateJob(companyId, request); httpError != nil {
		return httpError
	}
	util.WriteJSON(w, http.StatusOK, map[string]any{"message": "job created successfully"})
	return nil
}

func (jh *JobHandler) OfferJob(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	request := entity.OfferJobRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	if httpError = entity.ValidateOfferJobRequest(request); httpError != nil {
		return httpError
	}

	if httpError = jh.jobService.OfferJob(request); httpError != nil {
		return httpError
	}
	util.WriteJSON(w, http.StatusOK, map[string]any{"message": "job offered successfully"})
	return nil
}
