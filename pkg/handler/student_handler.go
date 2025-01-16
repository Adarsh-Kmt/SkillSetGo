package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/adarsh-kmt/skillsetgo/db"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/adarsh-kmt/skillsetgo/pkg/util"
	"github.com/gorilla/mux"
)

var dbConn *sql.DB

type StudentHandler struct {
	ss service.StudentService
}
type User struct {
	Name           string  `json:"name" binding:"required"`
	Branch         string  `json:"branch" binding:"required"`
	CGPA           float64 `json:"cgpa" binding:"required"`
	ActiveBacklogs bool    `json:"activebacklogs" binding:"required"`
	EmailID        string  `json:"emailid" binding:"required,email"`
	USN            string  `json:"usn" binding:"required"`
	CounsellorName string  `json:"counsellor" binding:"required"`
}

func NewStudentHandler(ss service.StudentService) *StudentHandler {
	return &StudentHandler{ss: ss}
}

func (sh *StudentHandler) MuxSetup(mux *mux.Router) *mux.Router {

	mux.HandleFunc("/register", util.MakeHttpHandlerFunc(sh.registerUser)).Methods("POST")
	return mux
}

func (sh *StudentHandler) registerUser(w http.ResponseWriter, r *http.Request) (httpError *util.HTTPError) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}
	arg := db.InsertUserParams{
		Usn:            user.USN,
		Name:           user.Name,
		Branch:         user.Branch,
		Cgpa:           user.CGPA,
		ActiveBacklogs: user.ActiveBacklogs,
		EmailID:        user.EmailID,
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
	return nil
}
