package main

import (
	"net/http"
	"project/handler"
	"project/router"
	"project/service"
	"project/structs"
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
	service := service.NewService()
	handler := handler.NewHandler(service)
	service.SetUsers(user)
	router := router.Init(handler)

	http.ListenAndServe(":80", router)

}
