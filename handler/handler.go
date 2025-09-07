package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/service"
	"project/structs"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service service.Service
}

// NewHandler конструктор для создания экземпляра Handler
// Принимает сервис через dependency injection
func NewHandler(service service.Service) *Handler {

	return &Handler{
		service: service,
	}

}

// GetList обработчик для получения списка всех пользователей
// GET /api/users
func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	list := h.service.Get()
	//Проверка длины строки
	if len(list) == 0 {
		fmt.Fprint(w, "Not found\n")
	}
	for key, value := range list {
		if value != nil {
			fmt.Fprintf(w, "%d: %s %s, Age: %d\n", key, value.Name, value.Lastname, value.Age)
		} else {

			fmt.Fprintf(w, "%d: [Delited user]\n", key)
		}
	}
}

// GetUser обработчик для получения конкретного пользователя по ID
// GET /api/users/:id
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%v %v %v\n", user.Name, user.Lastname, user.Age)
}

// AddUser обработчик для создания нового пользователя
// POST /api/users
func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Проверяем наличие ошибок
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var user structs.PostUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	if user.Name == "" || user.Lastname == "" {
		http.Error(w, "User name or lastname are empty", http.StatusBadRequest)
		return
	}

	//todo:получить данные из body
	id := h.service.Post(user.Name, user.Lastname, user.Age)
	fmt.Fprint(w, id)

}

// DeleteUser обработчик для удаления пользователя по ID
// DELETE /api/users/:id
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	err := h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "User deleted successfully")
}

// UpdateUser обработчик для обновления данных пользователя
// PUT /api/users/:id
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	name := r.URL.Query().Get("name")
	lastname := r.URL.Query().Get("lastname")
	age := r.URL.Query().Get("age")

	err := h.service.UpdateUser(id, name, lastname, age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "User update successfully")
}
