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
	companyService := service.NewCompanyServiceImpl()
	jobHandler := handler.NewJobHandler(jobService)
	studentHandler := handler.NewStudentHandler(studentService)
	companyHandler := handler.NewCompanyHandler(companyService)
	router := mux.NewRouter()
	router = jobHandler.MuxSetup(router)
	router = studentHandler.MuxSetup(router)
	router = companyHandler.MuxSetup(router)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
