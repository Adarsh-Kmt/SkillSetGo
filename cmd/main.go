package main

import (
	"log"
	"net/http"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/handler"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/gorilla/mux"
)

func main() {

	if err := db.PostgresDBClientInit(); err != nil {
		log.Fatal(err)
	}
	jobService := service.NewJobServiceImpl()
	studentService := service.NewStudentServiceImpl()
	authService := service.NewAuthServiceImpl()

	jobHandler := handler.NewJobHandler(jobService)
	studentHandler := handler.NewStudentHandler(studentService)
	authHandler := handler.NewAuthHandler(authService)

	router := mux.NewRouter()

	router = jobHandler.MuxSetup(router)
	router = studentHandler.MuxSetup(router)
	router = authHandler.MuxSetup(router)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
