package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, user models.UserData) error
}

type MailRepository struct {
	db *sql.DB

	conf    config.Postgres
	infoLog *logrus.Logger
}

func NewMailRepository(ctx context.Context, conf config.Postgres, infoLog *logrus.Logger) (*MailRepository, error) {
	url := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DB,
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// тайм-аут для подключения к бд
	timeoutCtx, cancel := context.WithTimeout(ctx, conf.Timeout)
	defer cancel()

	// пингуем бд, чтобы проверить, что она запущена и принимает соединения
	err = db.PingContext(timeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish connection with postgres: %v", err)
	}
	logrus.Debug("successfully established postgres connection!")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance: %v", err)
	}

	migrations := fmt.Sprintf("file://%v", conf.Migrations)
	m, err := migrate.NewWithDatabaseInstance(migrations, conf.DB, driver)
	if err != nil {
		return nil, fmt.Errorf("migrate.NewWithDatabaseInstance: %v", err)
	}

	// мигрируемся
	logrus.Debug("applying migrations...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Debug("nothing to migrate")
		} else {
			return nil, fmt.Errorf("error when migrating: %v", err)
		}
	} else {
		logrus.Debug("migrated successfully!")
	}

	return &MailRepository{
		db:      db,
		conf:    conf,
		infoLog: infoLog,
	}, nil
}

// создание пользователя в бд
func (mr *MailRepository) Create(ctx context.Context, user models.UserData) error {
	query := `
	INSERT INTO 
		email_schema.users (mail, nickname, firstname, lastname)
	VALUES 
		($1, $2, $3, $4);
	`

	if _, err := mr.db.ExecContext(
		ctx,
		query,
		user.Mail,
		user.Nickname,
		user.FirstName,
		user.LastName,
	); err != nil {
		if err.(*pq.Error).Code == "23505" {
			return ErrAlreadyExists
		}
		return fmt.Errorf("%v: %v", err.(*pq.Error).Code, err)
	}

	mr.infoLog.Info(fmt.Sprintf("created user: %+v", user))

	return nil
}
