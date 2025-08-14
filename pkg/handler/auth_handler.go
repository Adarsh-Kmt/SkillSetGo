package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/helper"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
func (ah *AuthHandler) MuxSetup(router *mux.Router) *mux.Router {

	router.HandleFunc("/student/login", helper.MakeHttpHandlerFunc(ah.LoginStudent)).Methods("POST")
	router.HandleFunc("/student/register", helper.MakeHttpHandlerFunc(ah.StudentRegister)).Methods("POST")
	router.HandleFunc("/company/login", helper.MakeHttpHandlerFunc(ah.CompanyLogin)).Methods("POST")
	router.HandleFunc("/company/register", helper.MakeHttpHandlerFunc(ah.CompanyRegister)).Methods("POST")

	return router
}

func (ah *AuthHandler) LoginStudent(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	var request entity.LoginStudentRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		slog.Error(err.Error())
		return &helper.HTTPError{StatusCode: 400, Error: "invalid request body"}
	}

	accessToken, httpError := ah.authService.LoginStudent(&request)
	if httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
	return nil
}

// StudentRegister TODO: implement me
func (ah *AuthHandler) StudentRegister(w http.ResponseWriter, r *http.Request) *helper.HTTPError {
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

	helper.WriteJSON(w, 200, map[string]string{"message": "student registered successfully"})
	return nil
}

func (ah *AuthHandler) CompanyLogin(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	var request entity.LoginCompanyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "invalid request body"}
	}

	accessToken, httpError := ah.authService.LoginCompany(&request)
	if httpError != nil {
		return httpError
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
	return nil
}

// CompanyRegister TODO: implement me
func (ah *AuthHandler) CompanyRegister(w http.ResponseWriter, r *http.Request) *helper.HTTPError {
	request := &entity.RegisterCompanyRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return &helper.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	httpError := entity.ValidateRegisterCompanyRequest(*request)

	if httpError != nil {
		return httpError
	}

	httpError = ah.authService.RegisterCompany(request)

	if httpError != nil {
		return httpError
	}

	helper.WriteJSON(w, 200, map[string]string{"message": "company registered successfully"})

	return nil
}
