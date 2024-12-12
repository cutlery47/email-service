package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/repo"
	"github.com/cutlery47/email-service/internal/service"
	"golang.org/x/exp/rand"
)

func Run() error {
	rand.Seed(uint64(time.Now().UnixNano()))
	ctx := context.Background()

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("error when setting up config: %v", err)
	}

	cache := service.NewMapCache(conf.Cache)
	repo := repo.NewMailRepository(ctx, conf.Postgres)

	service := service.NewMailService(cache, repo, conf.SMTP, conf.Code)

	log.Println(service)

	return nil
}
