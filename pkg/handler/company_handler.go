package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/gorilla/mux"
)

type CompanyHandler struct {
	companyService service.CompanyService
}

func NewCompanyHandler(cs service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: cs}
}

func (handler *CompanyHandler) MuxSetup(mux *mux.Router) *mux.Router {

	companyAdminRoleRequired := []string{"company admin"}
	mux.HandleFunc("/company/job", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.GetPublishedJobs), companyAdminRoleRequired)).Methods("GET")
	mux.HandleFunc("/company/job", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.CreateJob), companyAdminRoleRequired)).Methods("POST")
	mux.HandleFunc("/company/job/offer", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.OfferJob), companyAdminRoleRequired)).Methods("POST")
	mux.HandleFunc("/company/job/{job-id}/applicants", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.GetJobApplicants), companyAdminRoleRequired)).Methods("GET")
	mux.HandleFunc("/company/job/{job-id}/offer", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.GetOfferStatus), companyAdminRoleRequired)).Methods("GET")
	mux.HandleFunc("/company/job/interview", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.ScheduleInterview), companyAdminRoleRequired)).Methods("POST")

	mux.HandleFunc("/company/job/{job-id}/interview", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(handler.GetScheduledInterviews), companyAdminRoleRequired)).Methods("GET")

	mux.HandleFunc("/stats", helper.MakeHttpHandlerFunc(handler.GetPlacementStats)).Methods("GET")
	return mux
}

func (handler *CompanyHandler) CreateJob(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	companyId, err := helper.ValidateAccessToken(r.Header.Get("Auth"))

	if err != nil {
		return &helper.HTTPError{StatusCode: 500, Error: "internal server error"}
	}

	request := entity.CreateJobRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "bad request"}
	}
	//log.Println(request)
	if httpError = entity.ValidateCreateJobRequest(request); httpError != nil {
		return httpError
	}

	if httpError = handler.companyService.CreateJob(companyId, request); httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]any{"message": "job created successfully"})
	return nil
}

func (handler *CompanyHandler) OfferJob(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	request := entity.OfferJobRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	if httpError = entity.ValidateOfferJobRequest(request); httpError != nil {
		return httpError
	}

	if httpError = handler.companyService.OfferJob(request); httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]any{"message": "job offered successfully"})
	return nil
}

func (handler *CompanyHandler) GetPublishedJobs(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	companyId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))

	if httpError != nil {
		return httpError
	}

	jobs, httpError := handler.companyService.GetPublishedJobs(companyId)

	if httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
	return nil
}

func (handler *CompanyHandler) GetJobApplicants(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	companyId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))
	if httpError != nil {
		return httpError
	}

	vars := mux.Vars(r)

	jobIdString := vars["job-id"]

	jobId, err := strconv.Atoi(jobIdString)

	if err != nil || jobId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "invalid job id"}
	}

	profiles, httpError := handler.companyService.GetJobApplicants(companyId, jobId)

	if httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{"profiles": profiles})

	return nil
}

func (handler *CompanyHandler) GetOfferStatus(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	companyId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))
	if httpError != nil {
		return httpError
	}

	vars := mux.Vars(r)

	jobIdString := vars["job-id"]

	jobId, err := strconv.Atoi(jobIdString)

	if err != nil || jobId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "invalid job id"}
	}

	offers, httpError := handler.companyService.GetOfferStatus(companyId, jobId)
	if httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{"offers": offers})
	return nil
}

func (handler *CompanyHandler) ScheduleInterview(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	companyId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))
	if httpError != nil {
		return httpError
	}

	request := entity.ScheduleInterviewRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "bad request"}
	}
	httpError = entity.ValidateScheduleInterviewRequest(request)
	if httpError != nil {
		return httpError
	}

	httpError = handler.companyService.ScheduleInterview(companyId, request)
	if httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]any{"message": "interview scheduled successfully"})

	return nil
}

func (handler *CompanyHandler) GetScheduledInterviews(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	companyId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))
	if httpError != nil {
		return httpError
	}

	vars := mux.Vars(r)

	jobIdString := vars["job-id"]

	jobId, err := strconv.Atoi(jobIdString)

	if err != nil || jobId == 0 {
		return &helper.HTTPError{StatusCode: 400, Error: "invalid job id"}
	}

	interviews, httpError := handler.companyService.GetScheduledInterviews(companyId, jobId)
	if httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]any{"interviews": interviews})
	return nil
}

func (handler *CompanyHandler) GetPlacementStats(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	stats, httpError := handler.companyService.GetPlacementStats()

	if httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]any{"stats": stats})
	return nil
}
