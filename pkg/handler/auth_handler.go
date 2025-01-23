package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
func (ah *AuthHandler) MuxSetup(router *mux.Router) *mux.Router {

	router.HandleFunc("/student/login", util.MakeHttpHandlerFunc(ah.LoginStudent)).Methods("POST")
	router.HandleFunc("/student/register", util.MakeHttpHandlerFunc(ah.StudentRegister)).Methods("POST")
	router.HandleFunc("/company/login", util.MakeHttpHandlerFunc(ah.CompanyLogin)).Methods("POST")
	router.HandleFunc("/company/register", util.MakeHttpHandlerFunc(ah.CompanyRegister)).Methods("POST")

	return router
}

func (ah *AuthHandler) LoginStudent(w http.ResponseWriter, r *http.Request) *util.HTTPError {

	var request entity.LoginStudentRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "invalid request body"}
	}

	accessToken, httpError := ah.authService.LoginStudent(&request)
	if httpError != nil {
		return httpError
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
	return nil
}

// StudentRegister TODO: implement me
func (ah *AuthHandler) StudentRegister(w http.ResponseWriter, r *http.Request) *util.HTTPError {
	request := &entity.RegisterStudentRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}

	httpError := entity.ValidateRegisterStudentRequest(*request)

	if httpError != nil {
		return httpError
	}

	httpError = ah.authService.RegisterStudent(request)

	if httpError != nil {
		return httpError
	}

	util.WriteJSON(w, 200, map[string]string{"message": "student registered successfully"})
	return nil
}

func (ah *AuthHandler) CompanyLogin(w http.ResponseWriter, r *http.Request) *util.HTTPError {

	var request entity.LoginCompanyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "invalid request body"}
	}

	accessToken, httpError := ah.authService.LoginCompany(&request)
	if httpError != nil {
		return httpError
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
	return nil
}

// CompanyRegister TODO: implement me
func (ah *AuthHandler) CompanyRegister(w http.ResponseWriter, r *http.Request) *util.HTTPError {
	request := &entity.RegisterCompanyRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	httpError := entity.ValidateRegisterCompanyRequest(*request)

	if httpError != nil {
		return httpError
	}

	httpError = ah.authService.RegisterCompany(request)

	if httpError != nil {
		return httpError
	}

	util.WriteJSON(w, 200, map[string]string{"message": "company registered successfully"})

	return nil
}
