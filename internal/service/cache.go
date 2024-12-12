package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Кэш для хранения данных пользователей до момента подтверждения почты
type Cache interface {
	Put(user models.CachedUserData) error
	Get(mail string) (models.CachedUserData, error)
}

type MapCache struct {
	data map[string]models.CachedUserData
	mu   *sync.RWMutex

	conf    config.Cache
	infoLog *logrus.Logger
}

func NewMapCache(conf config.Cache, infoLog *logrus.Logger) *MapCache {
	cache := &MapCache{
		data: make(map[string]models.CachedUserData),
		mu:   &sync.RWMutex{},

		conf:    conf,
		infoLog: infoLog,
	}

	// запуск очистки кэша
	go cache.cleanup()
	return cache
}

func (mc *MapCache) Put(user models.CachedUserData) error {
	mc.mu.Lock()
	mc.data[user.Mail] = user
	mc.mu.Unlock()

	return nil
}

func (mc *MapCache) Get(mail string) (models.CachedUserData, error) {
	mc.mu.RLock()
	v, ok := mc.data[mail]
	mc.mu.RUnlock()

	if !ok {
		return models.CachedUserData{}, ErrCacheNotFound
	}

	return v, nil
}

func (mc *MapCache) cleanup() {
	for {
		time.Sleep(mc.conf.CleanupTimeout)
		mc.infoLog.Info("starting cache cleanup")

		mc.mu.Lock()
		for k, v := range mc.data {
			// проверка на превышение времени жизни ключа
			if v.CreatedAt.Add(v.ValidFor).Before(time.Now()) {
				delete(mc.data, k)
			}
		}
		mc.mu.Unlock()
	}
}

type RedisCache struct {
	cl *redis.Client

	conf config.Cache
}

func NewRedisCache(ctx context.Context, conf config.Cache) (*RedisCache, error) {
	url := fmt.Sprintf(
		"redis://%v:%v@%v:%v/%v",
		conf.Redis.Username,
		conf.Redis.Password,
		conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.DB,
	)

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	// пинг на проверку работоспособности
	pong := client.Ping(ctx)
	if pong.Err() != nil {
		return nil, pong.Err()
	}

	return &RedisCache{
		cl:   client,
		conf: conf,
	}, nil
}

func (rc *RedisCache) Put(user models.CachedUserData) error {
	return nil
}

func (rc *RedisCache) Get(mail string) (models.CachedUserData, error) {
	return models.CachedUserData{}, nil
}
