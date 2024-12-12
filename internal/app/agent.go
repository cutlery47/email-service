package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
	"github.com/cutlery47/email-service/internal/repo"
	"github.com/cutlery47/email-service/internal/service"
	"golang.org/x/exp/rand"
)

func RunAgent() error {
	rand.Seed(uint64(time.Now().UnixNano()))
	ctx := context.Background()

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("error when setting up config: %v", err)
	}

	// cache := service.NewMapCache(conf.Cache, logrus.StandardLogger())
	cache, err := service.NewRedisCache(ctx, conf.Cache)
	if err != nil {
		return fmt.Errorf("error when connecting to redis: %v", err)
	}
	repo := repo.NewMockRepository()
	service := service.NewMailService(cache, repo, conf.SMTP, conf.Code)

	err = service.Register(ctx, models.UserData{Mail: "ortizey47@gmail.com"})
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	code, _ := reader.ReadString('\n')

	err = service.Confirm(ctx, models.ConfirmationData{Mail: "ortizey47@gmail.com", Code: code[:len(code)-1]})
	if err != nil {
		return err
	}

	return nil
}
