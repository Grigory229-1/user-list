package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"project/structs"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

type TestService struct {
}

func (s *TestService) SetUsers(users []structs.User) {}
func (s *TestService) Get() []*structs.User {
	return []*structs.User{
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
	}
}
func (s *TestService) Post(name, lastname string, age uint) int {
	if name == "Alise" {
		return 0
	}
	return 1
}
func (s *TestService) Delete(idStr string) error {
	if idStr == "1" {
		return errors.New("user not found")
	}
	return nil
}
func (s *TestService) UpdateUser(idStr string, name, lastname string, age string) error {
	if idStr == "1" {
		return errors.New("user not found")
	}
	return nil
}
func (s *TestService) GetUser(idStr string) (*structs.User, error) {
	if idStr == "0" {
		return &structs.User{
			Name:     "Alise",
			Lastname: "Alise2",
			Age:      12,
		}, nil
	}
	return nil, errors.New("user not found")
}

func NewTestService() *TestService {
	return &TestService{}
}
func TestGetList(t *testing.T) {
	result := "0: Alise Alise2, Age: 12\n1: Max Alise2, Age: 0\n"
	h := NewHandler(NewTestService())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", nil)
	p := httprouter.Params{}
	h.GetList(w, r, p)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, result, w.Body.String())

}

// ////////////////////////////////////////////////////////
func TestGetUser(t *testing.T) {

	result := "Alise Alise2 12\n"
	h := NewHandler(NewTestService())
	w := httptest.NewRecorder()
	w1 := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/user/0", nil)
	p := httprouter.Params{httprouter.Param{Key: "id", Value: "0"}}
	p1 := httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}
	h.GetUser(w, r, p)
	h.GetUser(w1, r, p1)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, result, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w1.Code)
}

func TestAddUser(t *testing.T) {
	data := []byte(`{
    "name":"Alise",
    "lastname":"ALise2",
    "age":1
}`)
	body := bytes.NewReader(data)
	h := NewHandler(NewTestService())
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", body)
	p := httprouter.Params{}
	h.AddUser(w, r, p)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "0", w.Body.String())
}

func TestUppdateUser(t *testing.T) {
	h := NewHandler(NewTestService())
	w := httptest.NewRecorder()
	w1 := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/user/0?name=Max", nil)
	p := httprouter.Params{httprouter.Param{Key: "id", Value: "0"}}
	p1 := httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}
	h.UpdateUser(w, r, p)
	h.UpdateUser(w1, r, p1)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "User update successfully", w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w1.Code)

}

func TestDelitUser(t *testing.T) {

	h := NewHandler(NewTestService())
	w := httptest.NewRecorder()
	w1 := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/user/0", nil)
	p := httprouter.Params{httprouter.Param{Key: "id", Value: "0"}}
	p1 := httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}
	h.GetUser(w, r, p)
	h.GetUser(w1, r, p1)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
}

//////////////////////////////////////////////////////////
