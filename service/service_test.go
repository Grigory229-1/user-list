package service

import (
	"project/structs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteInvalidFormat(t *testing.T) {
	s := NewService()
	err := s.Delete("asddf")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid ID format")
}

func TestDeleteSuccess(t *testing.T) {
	s := NewService()
	s.SetUsers([]structs.User{
		{Name: "Tom", Lastname: "Max", Age: 20},
		{Name: "Alise", Lastname: "Max", Age: 20},
		{Name: "Max", Lastname: "Max", Age: 20},
	})
	err := s.Delete("0")
	assert.NoError(t, err)
	assert.Nil(t, s.user[0])
	assert.Equal(t, &structs.User{Name: "Alise", Lastname: "Max", Age: 20}, s.user[1])
	assert.Equal(t, &structs.User{Name: "Max", Lastname: "Max", Age: 20}, s.user[2])
}

func TestDeleteOutOfRange(t *testing.T) {
	s := NewService()
	err := s.Delete("0")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "user ID does not exist")

}

func TestUpdateUserInvalidFormat(t *testing.T) {
	s := NewService()
	err := s.UpdateUser("asdasf", "asdasf", "fffff", "90")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid ID format")
}

func TestUpdateUserSuccess(t *testing.T) {
	s := NewService()
	s.SetUsers([]structs.User{
		{Name: "Tom", Lastname: "Max", Age: 20},
		{Name: "Alise", Lastname: "Max", Age: 20},
		{Name: "Max", Lastname: "Max", Age: 20},
	})
	err := s.UpdateUser("1", "Max", "3w", "2")
	assert.NoError(t, err)
	assert.Equal(t, &structs.User{Name: "Tom", Lastname: "Max", Age: 20}, s.user[0])
	assert.Equal(t, &structs.User{Name: "Max", Lastname: "3w", Age: 2}, s.user[1])
	assert.Equal(t, &structs.User{Name: "Max", Lastname: "Max", Age: 20}, s.user[2])
}

func TestUpdateUserOutOfRange(t *testing.T) {
	s := NewService()
	err := s.UpdateUser("0", "Max", "dfdf", "0")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "user ID does not exist")
}

func TestGetUserInvalidFormat(t *testing.T) {
	s := NewService()
	_, err := s.GetUser("asdasf")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid ID format")
}

func TestGetUserSuccess(t *testing.T) {
	s := NewService()
	s.SetUsers([]structs.User{
		{Name: "Tom", Lastname: "Max", Age: 20},
		{Name: "Alise", Lastname: "Max", Age: 20},
		{Name: "Max", Lastname: "Max", Age: 20},
	})
	user, err := s.GetUser("1")
	assert.NoError(t, err)
	assert.Equal(t, &structs.User{Name: "Alise", Lastname: "Max", Age: 20}, user)

}

func TestUserIdDoesNotExist(t *testing.T) {
	s := NewService()

	_, err := s.GetUser("0")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "user ID does not exist")
}

func TestSetUsers(t *testing.T) {
	s := NewService()
	s.SetUsers([]structs.User{
		{Name: "Tom", Lastname: "Max", Age: 20},
		{Name: "Alise", Lastname: "Max", Age: 20},
		{Name: "Max", Lastname: "Max", Age: 20},
	})
	assert.Equal(t, &structs.User{Name: "Tom", Lastname: "Max", Age: 20}, s.user[0])
	assert.Equal(t, &structs.User{Name: "Alise", Lastname: "Max", Age: 20}, s.user[1])
	assert.Equal(t, &structs.User{Name: "Max", Lastname: "Max", Age: 20}, s.user[2])

}

func TestGet(t *testing.T) {
	s := NewService()

	s.SetUsers([]structs.User{
		{Name: "Tom", Lastname: "Max", Age: 20},
		{Name: "Alise", Lastname: "Max", Age: 20},
		{Name: "Max", Lastname: "Max", Age: 20},
	})
	users := s.Get()
	assert.Equal(t, &structs.User{Name: "Tom", Lastname: "Max", Age: 20}, users[0])
	assert.Equal(t, &structs.User{Name: "Alise", Lastname: "Max", Age: 20}, users[1])
	assert.Equal(t, &structs.User{Name: "Max", Lastname: "Max", Age: 20}, users[2])
}

func TestPost(t *testing.T) {
	s := NewService()
	s.SetUsers([]structs.User{
		{Name: "Tom", Lastname: "Max", Age: 20},
		{Name: "Alise", Lastname: "Max", Age: 20},
		{Name: "Max", Lastname: "Max", Age: 20},
	})
	id := s.Post("Nik", "dddd", 0)

	assert.Equal(t, 3, id)
	assert.Equal(t, &structs.User{Name: "Nik", Lastname: "dddd", Age: 0}, s.user[id])
}
