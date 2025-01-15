package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adarsh-Kmt/SkillSetGo/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Name           string  `json:"name" binding:"required"`
	Branch         string  `json:"branch" binding:"required"`
	CGPA           float64 `json:"cgpa" binding:"required"`
	ActiveBacklogs bool    `json:"activebacklogs" binding:"required"`
	EmailID        string  `json:"emailid" binding:"required,email"`
	USN            string  `json:"usn" binding:"required"`
	CounsellorName string  `json:"counsellor" binding:"required"`
}

var dbConn *sql.DB

func main() {
	var err error
	connStr := "user=postgres password=@neeshpostgres dbname=db sslmode=disable"
	dbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error: Failed to connect to database", err)
	}
	defer dbConn.Close()
	router := mux.NewRouter()
	router.HandleFunc("/register", registerUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User

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
