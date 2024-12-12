package models

import "time"

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
