package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/db"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
)

type StudentHandler struct {
	js service.JobService
}

func NewStudentHandler(js service.JobService) *StudentHandler {
	return &StudentHandler{js: js}
}

func (sh *StudentHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/register", util.MakeHttpHandlerFunc(sh.registerUser)).Methods("POST")
	return mux
}

func (sh *StudentHandler) registerUser(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}
	arg := db.InsertUserParams{
		Name:           user.Name,
		Branch:         user.Branch,
		Cgpa:           user.CGPA,
		ActiveBacklogs: user.ActiveBacklogs,
		EmailID:        user.EmailID,
		Usn:            user.USN,
		CounsellorName: user.CounsellorName,
	}
	queries := db.New(dbConn)
	err = queries.InsertUser(context.TODO(), arg)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json") //response and its type- json
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "successful"}); err != nil {
		panic(err)
	}
}
