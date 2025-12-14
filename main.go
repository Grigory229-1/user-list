package main

import (
	"log"
	"net/http"
	"project/handler"
	"project/psgservice"
	"project/router"
	"project/service"
	"project/structs"

	_ "github.com/lib/pq"
)

const (
	serviceProviderSlice = iota
	serviceProviderPostgresql
)

func main() {
	user := []structs.User{
		{
			Name:     "Alise",
			Lastname: "Alise2",
			Age:      12,
		},
		{
			Name:     "Max",
			Lastname: "Alise2",
			Age:      0,
		},
		{
			Name:     "ALex",
			Lastname: "Alise2",
			Age:      9,
		},
	}
	config := psgservice.Config{

		Host:     "localhost", // или имя/IP сервера (в Docker Compose часто "db")
		Port:     5432,
		User:     "postgres",
		Password: "1111",
		DBname:   "userlist",
		SSLmode:  "disable",
	}
	var serviceChoise service.Service
	//пеерключение serviceProvider смена с service на pgservice
	serviceProvider := serviceProviderSlice
	switch serviceProvider {
	case serviceProviderSlice:
		serviceChoise = service.NewService()
		serviceChoise.SetUsers(user)
	case serviceProviderPostgresql:
		serviceChoise = psgservice.NewService(config)
	}

	handler := handler.NewHandler(serviceChoise)
	router := router.Init(handler)
	log.Println("runnig The App on port :80")
	http.ListenAndServe(":80", router)

}
