package router

import (
	"project/handler"

	"github.com/julienschmidt/httprouter"
)

func Init(h *handler.Handler) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", h.GetList)
	router.POST("/", h.AddUser)
	router.GET("/:id", h.GetUser)
	router.DELETE("/:id", h.DeleteUser)
	router.PUT("/:id", h.UpdateUser)
	return router

}
