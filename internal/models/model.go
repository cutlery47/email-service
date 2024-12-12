package models

import "time"

// Данные, отправляемые пользователями при регистрации
// Эти данные заносятся в конечную бд
type UserData struct {
	Mail      string `json:"mail,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"secondname,omitempty"`
}

// Данные, отправляемые пользователями при подверждении регистрации
// По этим данным мы получаем информацию из временного хранилища
type ConfirmationData struct {
	Mail string `json:"mail,omitempty"`
	Code string `json:"code,omitempty"`
}

// Данные, хранящиеся во временном хранилище
type CachedUserData struct {
	UserData
	Code      string        `json:"code"`
	CreatedAt time.Time     `json:"created_at"`
	ValidFor  time.Duration `json:"valid_for"`
}
