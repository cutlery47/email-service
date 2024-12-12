package service

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
	"github.com/cutlery47/email-service/internal/repo"
	"golang.org/x/exp/rand"
)

type Service interface {
	Register(ctx context.Context, user models.UserData) error
	Confirm(ctx context.Context, data models.ConfirmationData) error
}

type MailService struct {
	cache Cache
	repo  repo.Repository
	auth  smtp.Auth

	codeConf config.Code
	smtpConf config.SMTP
}

func NewMailService(cache Cache, repo repo.Repository, smtpConf config.SMTP, codeConf config.Code) *MailService {
	auth := smtp.PlainAuth(
		"",
		smtpConf.Username,
		smtpConf.Password,
		smtpConf.Hostname,
	)

	return &MailService{
		cache: cache,
		repo:  repo,
		auth:  auth,

		codeConf: codeConf,
		smtpConf: smtpConf,
	}
}

// 1) получаем запрос на регистрацию
// 2) достаем из запроса email и отправляем юзеру код для подверждения
// 3) параллельно записываем передеанные данные во временнок хранилище (мапа, редис)
func (ms *MailService) Register(ctx context.Context, user models.UserData) error {
	code := ms.generateCode()

	cached := models.CachedUserData{
		UserData:  user,
		Code:      code,
		CreatedAt: time.Now(),
		ValidFor:  ms.codeConf.TTL,
	}

	if err := ms.cache.Put(cached); err != nil {
		return err
	}

	// jтсылаем пользователю сообщение с кодом подтверждения
	return smtp.SendMail(
		fmt.Sprintf("%v:%v", ms.smtpConf.Hostname, ms.smtpConf.Port),
		ms.auth,
		ms.smtpConf.Username,
		[]string{user.Mail},
		[]byte(fmt.Sprintf("Subject: Signup\nYour confirmation code %v", code)),
	)
}

// 1) получем код подтверждения
// 2) идем по переданной почте во временное хранилище и ищем там данные юзера,
// которые он отправил ранее
// 3) при совпадении кодов - делаем запись в бд, иначе отправляем ошибку
func (ms *MailService) Confirm(ctx context.Context, data models.ConfirmationData) error {
	cached, err := ms.cache.Get(data.Mail)
	if err != nil {
		return err
	}

	if data.Code != cached.Code {
		return ErrWrongCode
	}

	return ms.repo.Create(ctx, cached.UserData)
}

func (ms *MailService) generateCode() string {
	runes := []rune(ms.codeConf.Runes)
	runesLength := len(runes)

	code := make([]rune, ms.codeConf.Length)
	for i := range ms.codeConf.Length {
		code[i] = runes[rand.Intn(runesLength)]
	}

	return string(code)
}
