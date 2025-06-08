package service

import "fmt"

type Service interface {
	SetUsers(users []string) 
	Get() []string
	Post(name string) int
	Delete(id string) error
	UpdateUser(idStr string, newName string) error
	GetUser(idStr string) (string, error)
}

type ServiceImpl struct {
	user []string
}

func NewService() *ServiceImpl {
	user := make([]string, 0)
	return &ServiceImpl{
		user: user,
	}

}

func (s *ServiceImpl) SetUsers(users []string) {
	s.user = users
}

func (s *ServiceImpl) Get() []string {

	return s.user
}

func (s *ServiceImpl) Post(name string) int {
	s.user = append(s.user, name)
	return len(s.user) - 1

}

// Удалить пользователя по индексу

func (s *ServiceImpl) Delete(idStr string) error {
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return fmt.Errorf("invalid ID format :%v", err)
	}

	if id < 0 || id >= len(s.user) {
		return fmt.Errorf("user ID out of range")
	}

	// Удаляем пользователя по индексу
	s.user[id] = ""
	return nil
}

func (s *ServiceImpl) UpdateUser(idStr string, newName string) error {
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return fmt.Errorf("invalid ID format :%v", err)
	}
	if id >= len(s.user) {
		return fmt.Errorf("user ID out of range")
	}
	s.user[id] = newName
	return nil
}
func (s *ServiceImpl) GetUser(idStr string) (string, error) {
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return "", fmt.Errorf("invalid ID format :%v", err)
	}
	if id >= len(s.user) {
		return "", fmt.Errorf("user ID out of range")
	}
	return s.user[id], nil
}
