package handler

import (
	"encoding/json"
	"github.com/adarsh-kmt/skillsetgo/pkg/response"
	"net/http"
	"strconv"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
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

	router.HandleFunc("/student/offer", util.MakeAuthenticatedHandler(util.MakeHttpHandlerFunc(sh.GetJobOffers))).Methods("GET")
	router.HandleFunc("/student/apply/{job-id}", util.MakeAuthenticatedHandler(util.MakeHttpHandlerFunc(sh.ApplyForJob))).Methods("POST")
	router.HandleFunc("/student/offer", util.MakeAuthenticatedHandler(util.MakeHttpHandlerFunc(sh.GetJobOffers))).Methods("GET")
	router.HandleFunc("/student/offer", util.MakeAuthenticatedHandler(util.MakeHttpHandlerFunc(sh.PerformJobOfferAction))).Methods("PUT")
	return router
}

func (sh *StudentHandler) GetJobOffers(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	var (
		studentId int
		offers    []response.JobOfferResponse
	)

	if studentId, httpError = util.ValidateAccessToken(r.Header.Get("Auth")); httpError != nil {
		return httpError
	}

	if offers, httpError = sh.studentService.GetJobOffers(studentId); httpError != nil {
		return httpError
	}

	util.WriteJSON(w, 200, map[string]any{"offers": offers})

	return nil
}

func (sh *StudentHandler) ApplyForJob(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	var (
		studentId int
		jobId     int
		err       error
	)
	vars := mux.Vars(r)

	jobIdString := vars["job-id"]

	if jobId, err = strconv.Atoi(jobIdString); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "invalid job id"}
	}

	if httpError = sh.studentService.ApplyForJob(studentId, jobId); httpError != nil {
		return httpError
	}

	util.WriteJSON(w, 200, map[string]string{"response": "applied for job successfully"})
	return nil
}

func (sh *StudentHandler) PerformJobOfferAction(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	var (
		studentId int
		request   entity.PerformJobOfferActionRequest
		err       error
	)

	if studentId, httpError = util.ValidateAccessToken(r.Header.Get("Auth")); httpError != nil {
		return httpError
	}

	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	if httpError = entity.ValidatePerformJobOfferActionRequest(request); httpError != nil {
		return httpError
	}

	if httpError = sh.studentService.PerformJobOfferAction(studentId, request); httpError != nil {
		return httpError
	}

	util.WriteJSON(w, 200, map[string]string{"response": "action performed successfully"})

	return nil
}
