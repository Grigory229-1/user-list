package main

import (
	"net/http"
	"project/handler"
	"project/router"
	"project/service"
)

func main() {
	user := []string{"Tom", "Alise", "Max"}
	service := service.NewService()
	handler := handler.NewHandler(service)
	service.SetUsers(user)
	router := router.Init(handler)

	http.ListenAndServe(":80", router)
	
}
