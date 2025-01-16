package handler

import (
	"encoding/json"
	"net/http"

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

	mux.HandleFunc("/company", util.MakeHttpHandlerFunc(jh.GetCompanies)).Methods("GET")
	mux.HandleFunc("/company", util.MakeHttpHandlerFunc(jh.RegisterCompany)).Methods("POST")
	return mux
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
func (ch *CompanyHandler) RegisterCompany(w http.ResponseWriter, r *http.Request) error {
	var req service.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return util.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := ch.cs.RegisterCompany(r.Context(), req); err != nil {
		return util.NewHTTPError(http.StatusInternalServerError, "Failed to register company")
	}

	util.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Company registered successfully"})
	return nil
}

// GetCompanies handles fetching the list of companies
func (ch *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) error {
	companies, err := ch.cs.GetCompanies(r.Context())
	if err != nil {
		return util.NewHTTPError(http.StatusInternalServerError, "Failed to fetch companies")
	}

	util.RespondWithJSON(w, http.StatusOK, companies)
	return nil
}
