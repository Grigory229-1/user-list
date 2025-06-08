package service

import (
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
	s.SetUsers([]string{"Tom", "Alise", "Max"})
	err := s.Delete("0")
	assert.NoError(t, err)
	assert.Equal(t, "", s.user[0])
	assert.Equal(t, "Alise", s.user[1])
	assert.Equal(t, "Max", s.user[2])
}

func TestDeleteOutOfRange(t *testing.T) {
	s := NewService()
	err := s.Delete("0")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "user ID out of range")

}

func TestUpdateUserInvalidFormat(t *testing.T) {
	s := NewService()
	err := s.UpdateUser("asdasf", "asdasf")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid ID format")
}

func TestUpdateUserSuccess(t *testing.T) {
	s := NewService()
	s.SetUsers([]string{"Tom", "Alise", "Max"})
	err := s.UpdateUser("1", "Max")
	assert.NoError(t, err)
	assert.Equal(t, "Tom", s.user[0])
	assert.Equal(t, "Max", s.user[1])
	assert.Equal(t, "Max", s.user[2])
}

func TestUpdateUserOutOfRange(t *testing.T) {
	s := NewService()
	err := s.UpdateUser("0", "Max")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "user ID out of range")
}

func TestGetUserInvalidFormat(t *testing.T) {
	s := NewService()
	_, err := s.GetUser("asdasf")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid ID format")
}

func TestGetUserSuccess(t *testing.T) {
	s := NewService()
	s.SetUsers([]string{"Tom", "Alise", "Max"})
	user, err := s.GetUser("1")
	assert.NoError(t, err)
	assert.Equal(t, "Alise", user)

}

func TestGetUserOutOfRange(t *testing.T) {
	s := NewService()

	_, err := s.GetUser("0")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "user ID out of range")
}

func TestSetUsers(t *testing.T) {
	s := NewService()
	s.SetUsers([]string{"Tom", "Alise", "Max"})
	assert.Equal(t, "Tom", s.user[0])
	assert.Equal(t, "Alise", s.user[1])
	assert.Equal(t, "Max", s.user[2])

}

func TestGet(t *testing.T) {
	s := NewService()
	
	s.SetUsers([]string{"Tom", "Alise", "Max"})
	users := s.Get()
	assert.Equal(t, "Tom", users[0])
	assert.Equal(t, "Alise", users[1])
	assert.Equal(t, "Max", users[2])
}

func TestPost(t *testing.T) {
	s := NewService()
	s.SetUsers([]string{"Tom", "Alise", "Max"})
	id := s.Post("Nik")

	assert.Equal(t, 3, id)
	assert.Equal(t, "Nik", s.user[id])
}
