package models

import "time"

// Данные, отправляемые пользователями при регистрации
// Эти данные заносятся в конечную бд
type UserData struct {
	Mail      string
	Nickname  string
	FirstName string
	LastName  string
}

// Данные, отправляемые пользователями при подверждении регистрации
// По этим данным мы получаем информацию из временного хранилища
type ConfirmationData struct {
	Mail string
	Code string
}

// Данные, хранящиеся во временном хранилище
type CachedUserData struct {
	UserData
	Code      string
	CreatedAt time.Time
	ValidFor  time.Duration
}
