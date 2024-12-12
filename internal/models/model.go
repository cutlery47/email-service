package models

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
