package router

import (
	"project/handler"

	"github.com/julienschmidt/httprouter"
)

func Init(h *handler.Handler) *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/users/", h.GetList)
	router.POST("/api/users/", h.AddUser)
	router.GET("/api/users/:id", h.GetUser)
	router.DELETE("/api/users/:id", h.DeleteUser)
	router.PUT("/api/users/:id", h.UpdateUser)
	return router

}
