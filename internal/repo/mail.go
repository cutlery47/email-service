package repo

import (
	"context"
	"database/sql"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
)

type Repository interface {
	Create(ctx context.Context, user models.UserData) error
}

type MailRepository struct {
	db *sql.DB

	conf config.Postgres
}

func NewMailRepository(ctx context.Context, conf config.Postgres) *MailRepository {
	return &MailRepository{}
}

func (mr *MailRepository) Create(ctx context.Context, user models.UserData) error {
	return nil
}
