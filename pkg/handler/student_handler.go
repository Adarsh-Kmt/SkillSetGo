package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/pkg/entity"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
)

var dbConn *sql.DB

type StudentHandler struct {
	ss service.StudentService
}

func NewStudentHandler(ss service.StudentService) *StudentHandler {
	return &StudentHandler{ss: ss}
}

func (sh *StudentHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/register", util.MakeHttpHandlerFunc(sh.registerUser)).Methods("POST")
	return mux
}

func (sh *StudentHandler) registerUser(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {

	request := &entity.RegisterStudentRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}
	// arg := db.InsertUserParams{
	// 	Usn:            user.USN,
	// 	Name:           user.Name,
	// 	Branch:         user.Branch,
	// 	Cgpa:           user.CGPA,
	// 	ActiveBacklogs: user.ActiveBacklogs,
	// 	EmailID:        user.EmailID,
	// }

	// queries := db.New(dbConn)
	// err = queries.InsertUser(context.TODO(), arg)
	// if err != nil {
	// 	http.Error(w, "Failed to register user", http.StatusInternalServerError)
	// 	return
	// }
	if httpError = entity.ValidateRegisterStudentRequest(*request); httpError != nil {
		return httpError
	}
	if httpError = sh.ss.RegisterStudent(*request); httpError != nil {
		return httpError
	}
	w.Header().Add("Content-type", "application/json") //response and its type- json
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "successful"}); err != nil {
		panic(err)
	}
	return nil
}
