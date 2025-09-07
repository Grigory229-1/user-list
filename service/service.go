package service

import (
	"fmt"
	"project/structs"
)

// Service интерфейс определяет контракт для сервиса работы с пользователями
// Интерфейс позволяет абстрагироваться от конкретной реализации
// и легко подменять реализацию для тестирования
type Service interface {
	SetUsers(users []structs.User)
	Get() []*structs.User
	Post(name, lastname string, age uint) int
	Delete(idStr string) error
	UpdateUser(idStr string, name, lastname string, age string) error
	GetUser(idStr string) (*structs.User, error)
}

// ServiceImpl - реализация интерфейса Service
// Содержит slice указателей на User для эффективного управления памятью
type ServiceImpl struct {
	user []*structs.User
}

// NewService - конструктор для создания нового экземпляра ServiceImpl
// Инициализирует пустой slice пользователей
// Возвращает указатель на ServiceImpl для избежания копирования
func NewService() *ServiceImpl {
	user := make([]*structs.User, 0)
	return &ServiceImpl{
		user: user,
	}

}

// SetUsers добавляет массив пользователей в сервис
// Принимает slice структур User, сохраняет указатели на них
// Это позволяет работать с оригинальными объектами, а не копиями
func (s *ServiceImpl) SetUsers(users []structs.User) {
	for i := 0; i < len(users); i++ {
		s.user = append(s.user, &users[i])
	}
}

// Get возвращает всех пользователей в виде slice указателей
// Позволяет получать доступ к оригинальным объектам пользователей
// Возвращаются указатели для избежания копирования больших структур
func (s *ServiceImpl) Get() []*structs.User {

	return s.user
}

// Post создает нового пользователя и добавляет его в коллекцию
// Принимает основные данные пользователя: имя, фамилию, возраст
// Возвращает индекс созданного пользователя в slice
// Это упрощенная реализация - в production лучше использовать UUID
func (s *ServiceImpl) Post(name, lastname string, age uint) int {
	user := &structs.User{
		Name:     name,
		Lastname: lastname,
		Age:      age,
	}
	s.user = append(s.user, user)
	return len(s.user) - 1

}

// Delete удаляет пользователя по его ID (индексу в slice)
// Принимает ID в виде строки, парсит его в int
// Выполняет валидацию ID перед удалением
// Устанавливает значение nil вместо удаления из slice для сохранения индексов
// Удалить пользователя по индексу

func (s *ServiceImpl) Delete(idStr string) error {
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return fmt.Errorf("invalid ID format :%v", err)
	}

	if id < 0 || id >= len(s.user) || s.user[id] == nil {
		return fmt.Errorf("user ID does not exist")
	}

	// Удаляем пользователя по индексу
	s.user[id] = nil
	return nil
}

// UpdateUser обновляет данные существующего пользователя
// Принимает ID пользователя и новые значения (если не пустые)
// Поддерживает частичное обновление - только изменяемые поля
// Выполняет валидацию перед обновлением
func (s *ServiceImpl) UpdateUser(idStr string, name, lastname string, age string) error {
	id := 0
	var ageUint uint
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return fmt.Errorf("invalid ID format :%v", err)
	}
	if id >= len(s.user) || s.user[id] == nil {
		return fmt.Errorf("user ID does not exist")
	}

	if name != "" {
		s.user[id].Name = name
	}
	if lastname != "" {
		s.user[id].Lastname = lastname
	}
	if age != "" {
		_, err := fmt.Sscanf(age, "%d", &ageUint)
		if err != nil {
			fmt.Printf("invalid ID format :%v\n", err)
		} else {
			s.user[id].Age = ageUint
		}
	}
	return nil
}

// GetUser возвращает пользователя по его ID
// Принимает строковый ID, парсит его и выполняет валидацию
// Возвращает указатель на пользователя или ошибку если пользователь не найден
// /////////////////////////////////////////////////////////////////////
func (s *ServiceImpl) GetUser(idStr string) (*structs.User, error) {
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format :%v", err)
	}
	if id >= len(s.user) || s.user[id] == nil {
		return nil, fmt.Errorf("user ID does not exist")
	}
	return s.user[id], nil
}
