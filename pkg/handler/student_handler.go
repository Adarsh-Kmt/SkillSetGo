package handler

import (
	"encoding/json"
	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"net/http"
	"strconv"

	"github.com/adarsh-kmt/skillsetgo/pkg/response"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/gorilla/mux"
)

type StudentHandler struct {
	studentService service.StudentService
}

func NewStudentHandler(ss service.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: ss,
	}
}

func (sh *StudentHandler) MuxSetup(router *mux.Router) *mux.Router {

	studentRoleRequired := []string{"student"}
	router.HandleFunc("/student/job/offer", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(sh.GetJobOffers), studentRoleRequired)).Methods("GET")
	router.HandleFunc("/student/job/{job-id}/apply", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(sh.ApplyForJob), studentRoleRequired)).Methods("POST")
	router.HandleFunc("/student/job/offer", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(sh.PerformJobOfferAction), studentRoleRequired)).Methods("PUT")
	router.HandleFunc("/student/job", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(sh.GetJobs), studentRoleRequired)).Methods("GET")
	router.HandleFunc("/student/{student-id}/profile", helper.MakeAuthenticatedHandler(helper.MakeHttpHandlerFunc(sh.GetStudentProfile))).Methods("GET")
	router.HandleFunc("/student/job/apply", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(sh.GetAlreadyAppliedJobs), studentRoleRequired)).Methods("GET")

	router.HandleFunc("/student/job/interview", helper.MakeAuthorizedHandler(helper.MakeHttpHandlerFunc(sh.GetScheduledInterviews), studentRoleRequired)).Methods("GET")
	return router
}

func (sh *StudentHandler) GetJobs(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	var (
		jobs   []*sqlc.GetJobsRow
		userId int
	)

	if userId, httpError = helper.ValidateAccessToken(r.Header.Get("Auth")); httpError != nil {
		return httpError
	}
	queryParams := r.URL.Query()
	salaryTierList := queryParams["salary-tier"]
	companyList := queryParams["company"]
	jobRoleList := queryParams["job-role"]

	for _, salaryTier := range salaryTierList {
		if salaryTier != "Dream" && salaryTier != "Open Dream" && salaryTier != "Mass Recruitment" {
			return &helper.HTTPError{StatusCode: 400, Error: "invalid salary tier url query parameter"}
		}
	}

	//log.Println(salaryTierList)
	if jobs, httpError = sh.studentService.GetJobs(userId, salaryTierList, jobRoleList, companyList); httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
	return nil
}
func (sh *StudentHandler) GetJobOffers(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	var (
		studentId int
		offers    []response.JobOfferResponse
	)

	if studentId, httpError = helper.ValidateAccessToken(r.Header.Get("Auth")); httpError != nil {
		return httpError
	}

	if offers, httpError = sh.studentService.GetJobOffers(studentId); httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, 200, map[string]any{"offers": offers})

	return nil
}

func (sh *StudentHandler) ApplyForJob(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	var (
		studentId int
		jobId     int
		err       error
	)

	if studentId, httpError = helper.ValidateAccessToken(r.Header.Get("Auth")); httpError != nil {
		return httpError
	}
	vars := mux.Vars(r)

	jobIdString := vars["job-id"]

	if jobId, err = strconv.Atoi(jobIdString); err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "invalid job id"}
	}

	if httpError = sh.studentService.ApplyForJob(studentId, jobId); httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, 200, map[string]string{"response": "applied for job successfully"})
	return nil
}

func (sh *StudentHandler) PerformJobOfferAction(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	var (
		studentId int
		request   entity.PerformJobOfferActionRequest
		err       error
	)

	if studentId, httpError = helper.ValidateAccessToken(r.Header.Get("Auth")); httpError != nil {
		return httpError
	}

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	if httpError = entity.ValidatePerformJobOfferActionRequest(request); httpError != nil {
		return httpError
	}

	if httpError = sh.studentService.PerformJobOfferAction(studentId, request); httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, 200, map[string]string{"response": "action performed successfully"})

	return nil
}

func (sh *StudentHandler) GetStudentProfile(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {
    vars := mux.Vars(r)
    studentIdStr := vars["student-id"]

    if studentIdStr == "" {
        return &helper.HTTPError{StatusCode: 400, Error: "invalid student id"}
    }

    studentId, err := strconv.Atoi(studentIdStr)
    if err != nil {
        return &helper.HTTPError{StatusCode: 400, Error: "invalid student id format"}
    }

    profile, httpError := sh.studentService.GetStudentProfile(studentId)
    if httpError != nil {
        return httpError
    }

    helper.WriteJSON(w, 200, map[string]any{"profile": profile})
    return nil
}

func (sh *StudentHandler) GetAlreadyAppliedJobs(w http.ResponseWriter, r *http.Request) (httpError *helper.HTTPError) {

	studentId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))
	if httpError != nil {
		return httpError
	}

	appliedJobs, httpError := sh.studentService.GetAlreadyAppliedJobs(studentId)
	if httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, 200, map[string]any{"jobs": appliedJobs})

	return nil
}

func (sh *StudentHandler) GetScheduledInterviews(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	studentId, httpError := helper.ValidateAccessToken(r.Header.Get("Auth"))
	if httpError != nil {
		return httpError
	}

	interviews, httpError := sh.studentService.GetScheduledInterviews(studentId)

	if httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, 200, map[string]any{"interviews": interviews})
	return nil
}
