package handler

import (
	"fmt"
	"net/http"
	"project/service"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {

	return &Handler{
		service: service,
	}

}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	list := h.service.Get()
	for key, value := range list {
		fmt.Fprintf(w, "%v %s\n", key, value)
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s\n", user)
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.URL.Query().Get("name")
	id := h.service.Post(name)
	fmt.Fprint(w, id)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	err := h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "User deleted successfully")
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	name := r.URL.Query().Get("name")
	err := h.service.UpdateUser(id, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "User update successfully")
}
