// Пакет psgservice отвечает за доступ к БД (PostgreSQL) через database/sql.
// ВНИМАНИЕ: текущая реализация открывает соединение к БД в каждом методе.
// В продакшене обычно хранят *sql.DB (пул соединений) в Service и переиспользуют его.
package psgservice

import (
	"database/sql"
	"fmt"
	"project/structs"
	"strconv"
)

// Service — «тонкий» сервис над БД. Сейчас хранит только строку подключения (DSN).
type Service struct {
	connStr string // DSN для postgres, формата: "host=... port=... user=... password=... dbname=... sslmode=..."
}

// Config — параметры подключения. Значения обычно совпадают с теми, что вы вводите в pgAdmin.
type Config struct {
	Host     string // хост (в Docker Compose часто имя сервиса: "db"; локально — "localhost")
	Port     int    // порт PostgreSQL (по умолчанию 5432)
	User     string // пользователь БД (например, "postgres" или ваш app-юзер)
	Password string // пароль
	DBname   string // имя базы данных
	SSLmode  string // режим SSL: "disable"/"require"/"verify-full" и т.д.
}

// NewService — конструктор сервиса. Собирает connStr (DSN) из config и возвращает Service.
// Обратите внимание: сам *sql.DB здесь не создаётся — он открывается внутри каждого метода.
func NewService(config Config) *Service {
	// Формируем DSN в key=value синтаксисе (поддерживается драйвером lib/pq).
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBname, config.SSLmode,
	)
	return &Service{connStr: connStr}
}

// SetUsers — заготовка под массовую запись пользователей.
// Сейчас пустой (ничего не делает). Можно реализовать bulk insert/update.
func (s *Service) SetUsers(users []structs.User) {
	// TODO: реализовать массовую загрузку при необходимости.
}

// Get — получение всех пользователей из таблицы userlist.
// ВАЖНО: используется "select *" и нет rows.Close()/rows.Err(), это потенциальные грабли.
// Также порядок колонок при SELECT * может не совпасть с порядком в Scan.
func (s *Service) Get() ([]*structs.User, error) {
	// Открываем соединение (создаёт/берёт из пула новый *sql.DB).
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		// В случае ошибки возвращаем пустой слайс и саму ошибку.
		return []*structs.User{}, err
	}
	defer db.Close() // Закрываем при выходе из метода (для *sql.DB это освобождение ресурсов, сами коннекты пулит драйвер).

	// Исправлено: используем явный список полей и ORDER BY для гарантированного порядка.
	rows, err := db.Query("SELECT id, user_name, lastname, age FROM userlist ORDER BY id")
	if err != nil {
		return []*structs.User{}, err
	}
	defer rows.Close() // Закрываем rows для освобождения ресурсов.

	userList := make([]*structs.User, 0)

	for rows.Next() {
		user := &structs.User{}
		// Исправлено: порядок колонок соответствует SELECT (id, user_name, lastname, age).
		var ageInt int
		err := rows.Scan(&user.Id, &user.Name, &user.Lastname, &ageInt)
		user.Age = uint(ageInt)
		if err != nil {
			// На ошибке логируем и продолжаем.
			fmt.Println(err)
			continue
		}
		userList = append(userList, user)
	}

	// Проверяем ошибки итератора после цикла.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}

// Post — создаёт пользователя. Использует INSERT ... RETURNING id для получения ID созданной записи.
func (s *Service) Post(name, lastname string, age uint) (int, error) {
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Исправлено: один INSERT с RETURNING id для получения ID созданной записи.
	var id int
	err = db.QueryRow(
		"INSERT INTO userlist (user_name, lastname, age) VALUES ($1, $2, $3) RETURNING id",
		name, lastname, age,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Delete — удаляет пользователя по id.
// ВАЖНО: метод предполагает наличие поля id в таблице userlist.
// Если id отсутствует в схеме — запрос завершится ошибкой.
func (s *Service) Delete(idStr string) error {
	// Преобразуем id из строки в int, валидируем.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %q: %w", idStr, err)
	}

	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Удаляем строку по id.
	// Хорошая практика — проверить RowsAffected() и вернуть sql.ErrNoRows, если ничего не удалено.
	_, err = db.Exec("delete from userlist where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser — обновляет пользователя по id.
// ВАЖНО: здесь есть критическая ошибка: Scan вызывается без & (указателей).
// Также age приходит строкой — его нужно парсить/приводить к числу перед UPDATE.
func (s *Service) UpdateUser(idStr string, name, lastname string, age string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %q: %w", idStr, err)
	}

	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Считываем текущие значения пользователя (по id).
	row := db.QueryRow("SELECT user_name, lastname, age FROM userlist WHERE id = $1", id)
	user := &structs.User{}

	// Исправлено: Scan теперь принимает указатели на поля.
	var ageInt int
	err = row.Scan(&user.Name, &user.Lastname, &ageInt)
	user.Age = uint(ageInt)
	if err != nil {
		return err
	}

	// Если какие-то поля пришли пустыми — оставляем прежние значения.
	if name == "" {
		name = user.Name
	}
	if lastname == "" {
		lastname = user.Lastname
	}
	if age == "" {
		age = fmt.Sprint(user.Age) // строковое представление возраста из структуры
	}

	// Обновляем запись. Здесь age подставляется как строка — БД приведёт тип автоматически,
	// но безопаснее заранее сконвертировать age в int.
	_, err = db.Exec(
		"update userlist set user_name = $1, lastname = $2, age = $3 where id = $4",
		name, lastname, age, id,
	)
	if err != nil {
		// ПЛОХО: panic в библиотечном коде. Лучше вернуть err наверх.
		return err
	}
	return nil
}

// GetUser — возвращает одного пользователя по id.
// Здесь запрос выбирает явный список колонок и корректно сканирует с & — это хорошо.
func (s *Service) GetUser(idStr string) (*structs.User, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id %q: %w", idStr, err)
	}

	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Выбираем явные поля — правильно (устойчиво к изменениям схемы).
	row := db.QueryRow("SELECT id, user_name, lastname, age FROM userlist WHERE id = $1", id)

	user := &structs.User{}
	// Исправлено: сканируем в &user.Id, а не в &id (который является параметром функции).
	var ageInt int
	err = row.Scan(&user.Id, &user.Name, &user.Lastname, &ageInt)
	user.Age = uint(ageInt)

	if err != nil {

		// Если запись не найдена, err будет sql.ErrNoRows — имеет смысл пробрасывать его как есть.
		return nil, err
	}

	return user, nil
}
