package service

import (
	"github.com/cutlery47/email-service/internal/models"
	"github.com/cutlery47/email-service/internal/repo"
)

type Service interface {
	Register(user models.UserData) error
	Confirm(data models.ConfirmationData) error
}

type MailService struct {
	cache Cache
	repo  repo.Repository
}

// 1) получаем запрос на регистрацию
// 2) достаем из запроса email и отправляем юзеру код для подверждения
// 3) параллельно записываем передеанные данные во временнок хранилище (мапа, редис)

// Предусмотреть
// 1) время жизни кода, время создания кода в бд (отдельные поля)
func (ms *MailService) Register(user models.UserData) error {

	return nil
}

// 1) получем код подтверждения
// 2) идем по переданной почте во временное хранилище и ищем там данные юзера,
// которые он отправил ранее
// 3) при совпадении кодов - делаем запись в бд, иначе отправляем ошибку
func (ms *MailService) Confirm(data models.ConfirmationData) error {
	return nil
}
