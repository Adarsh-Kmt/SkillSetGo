package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
)

type CompanyHandler struct {
	cs service.CompanyService
}

func NewCompanyHandler(cs service.CompanyService) *CompanyHandler {
	return &CompanyHandler{cs: cs}
}

func (ch *CompanyHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/company", util.MakeHttpHandlerFunc(ch.GetCompanies)).Methods("GET")
	mux.HandleFunc("/company", util.MakeHttpHandlerFunc(ch.RegisterCompany)).Methods("POST")
	return mux
}

// RegisterCompany handles the registration of a new company
func (ch *CompanyHandler) RegisterCompany(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {
	req := &entity.RegisterCompanyRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}
	// arg:=sqlc.CreateCompanyParams{
	// 	CompanyName: req.CompanyName,
	// 	PocName: req.PocName,
	// 	PocPhno: req.PocPhno,
	// 	Industry:req.Industry,
	// }
	if httpError := entity.ValidateRegisterCompanyRequest(*req); err != nil {
		return httpError
	}

	if httpError := ch.cs.RegisterCompany(*req); httpError != nil {
		return httpError
	}
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Company registered successfully"})
	return nil
}

// GetCompanies handles fetching the list of companies
func (ch *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) *util.HTTPError {
	request := &entity.RegisterCompanyRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return &util.HTTPError{StatusCode: 400, Error: "bad request"}
	}

	if httpError := entity.ValidateRegisterCompanyRequest(*request); httpError != nil {
		return httpError
	}

	return nil
}
