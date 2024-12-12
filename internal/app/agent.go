package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
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

	d := models.UserData{Mail: "ortizey47@gmail.com"}
	if err := cache.Put(ctx, models.CachedUserDataIn{
		CachedUserDataOut: models.CachedUserDataOut{
			UserData: d,
			Code:     "123123",
		},
		CreatedAt: time.Now(),
		ValidFor:  time.Hour,
	}); err != nil {
		log.Println(err)
	}

	c, err := cache.Get(ctx, "ortizey47@gmail.com")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(c)

	// repo := repo.NewMockRepository()
	// service := service.NewMailService(cache, repo, conf.SMTP, conf.Code)

	// err = service.Register(ctx, models.UserData{Mail: "ortizey47@gmail.com"})
	// if err != nil {
	// 	return err
	// }

	// reader := bufio.NewReader(os.Stdin)
	// code, _ := reader.ReadString('\n')

	// err = service.Confirm(ctx, models.ConfirmationData{Mail: "ortizey47@gmail.com", Code: code[:len(code)-1]})
	// if err != nil {
	// 	return err
	// }

	return nil
}
