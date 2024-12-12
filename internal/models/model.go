package models

import "time"

// Данные, отправляемые пользователями при регистрации
// Эти данные заносятся в конечную бд
type UserData struct {
	Mail      string `json:"mail,omitempty" redis:"mail"`
	Nickname  string `json:"nickname,omitempty" redis:"nickname"`
	FirstName string `json:"firstname,omitempty" redis:"firstname"`
	LastName  string `json:"secondname,omitempty" redis:"lastname"`
}

// Данные, отправляемые пользователями при подверждении регистрации
// По этим данным мы получаем информацию из временного хранилища
type ConfirmationData struct {
	Mail string `json:"mail,omitempty"`
	Code string `json:"code,omitempty"`
}

// Данные, получаемые из временного хранилища
type CachedUserDataOut struct {
	UserData
	Code string `json:"code" redis:"code"`
}

// Данные, записываемые во временное хранилище
type CachedUserDataIn struct {
	CachedUserDataOut
	CreatedAt time.Time     `json:"created_at" redis:"created_at"`
	ValidFor  time.Duration `json:"valid_for"`
}
