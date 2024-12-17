package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
	"github.com/sirupsen/logrus"
)

var (
	conf = config.Cache{
		CacheType:      "map",
		CleanupTimeout: 1 * time.Minute,
	}

	user = models.CachedUserDataIn{
		CachedUserDataOut: models.CachedUserDataOut{
			UserData: models.UserData{},
		},
		CreatedAt: time.Now(),
		ValidFor:  time.Second * 15,
	}

	cache Cache
	log   *logrus.Logger
	ctx   context.Context
)

func setup() {
	ctx = context.Background()
	log = logrus.StandardLogger()
	cache = NewMapCache(conf, log)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestPut(t *testing.T) {
	err := cache.Put(ctx, user)
	if err != nil {
		t.Fatal("error: ", err)
	}
}

func TestGet(t *testing.T) {
	data, err := cache.Get(ctx, user.Mail)
	if err != nil {
		t.Fatal("error: ", err)
	}

	if data != user.CachedUserDataOut {
		t.Fatal("error: data does not match")
	}
}
