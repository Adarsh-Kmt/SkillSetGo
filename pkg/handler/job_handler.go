package handler

import (
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type JobHandler struct {
	js service.JobService
}

func NewJobHandler(js service.JobService) *JobHandler {
	return &JobHandler{js: js}
}
func (jh *JobHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/job", util.MakeHttpHandlerFunc(jh.GetJobs)).Methods("GET")
	return mux
}

func (jh *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	var (
		jobs []*sqlc.GetJobsRow
	)

	queryParams := r.URL.Query()
	salaryTierList := queryParams["salary-tier"]

	if salaryTierList != nil {
		for _, salaryTier := range salaryTierList {
			if salaryTier != "Dream" && salaryTier != "Open Dream" {
				return &util.HTTPError{StatusCode: 400, Error: "invalid salary tier url query parameter."}
			}
		}
	}
	log.Println(salaryTierList)
	if jobs, httpError = jh.js.GetJobs(1, salaryTierList); httpError != nil {
		return httpError
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
	return nil
}
