package main

import (
	"database/sql"
	"log"
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

var db *sql.DB

func main() {
	var err error
	connStr := "user=user password=password dbname=db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error: Failed to connect to database", err)
	}
	defer db.Close()
}
