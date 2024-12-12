package service

import (
	"context"
	"encoding/json"
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
	Put(ctx context.Context, user models.CachedUserDataIn) error
	Get(ctx context.Context, mail string) (models.CachedUserDataOut, error)
}

type MapCache struct {
	data map[string]models.CachedUserDataIn
	mu   *sync.RWMutex

	conf    config.Cache
	infoLog *logrus.Logger
}

func NewMapCache(conf config.Cache, infoLog *logrus.Logger) *MapCache {
	cache := &MapCache{
		data: make(map[string]models.CachedUserDataIn),
		mu:   &sync.RWMutex{},

		conf:    conf,
		infoLog: infoLog,
	}

	// запуск очистки кэша
	go cache.cleanup()
	return cache
}

func (mc *MapCache) Put(ctx context.Context, user models.CachedUserDataIn) error {
	mc.mu.Lock()
	mc.data[user.Mail] = user
	mc.mu.Unlock()

	return nil
}

func (mc *MapCache) Get(ctx context.Context, mail string) (models.CachedUserDataOut, error) {
	mc.mu.RLock()
	v, ok := mc.data[mail]
	mc.mu.RUnlock()

	if !ok {
		return models.CachedUserDataOut{}, ErrCacheNotFound
	}

	return models.CachedUserDataOut{
		UserData: v.UserData,
		Code:     v.Code,
	}, nil
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

func (rc *RedisCache) Put(ctx context.Context, user models.CachedUserDataIn) error {
	var redisUser map[string]interface{}

	// конвертируем CachedUserDataIn в дефолтную мапу, иначе редис не пишет ключи
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v:", err)
	}

	if err := json.Unmarshal(data, &redisUser); err != nil {
		return fmt.Errorf("json.Unmarshal: %v:", err)
	}

	// запись в редис
	if cmd := rc.cl.HSet(ctx, user.Mail, redisUser); cmd.Err() != nil {
		return fmt.Errorf("rc.cl.HSet: %v", cmd.Err())
	}

	// установка ttl
	if cmd := rc.cl.Expire(ctx, user.Mail, user.ValidFor); cmd.Err() != nil {
		return fmt.Errorf("rc.cl.Expire: %v", cmd.Err())
	}

	return nil
}

func (rc *RedisCache) Get(ctx context.Context, mail string) (models.CachedUserDataOut, error) {
	var redisUser models.CachedUserDataOut

	// получаем значения по ключу
	mapResult, err := rc.cl.HGetAll(ctx, mail).Result()
	if err != nil {
		return models.CachedUserDataOut{}, fmt.Errorf("rc.cl.HGetAll: %v", err)
	}

	// маршалим в удобную структуру
	data, err := json.Marshal(mapResult)
	if err != nil {
		return models.CachedUserDataOut{}, fmt.Errorf("json.Marshal: %v", err)
	}

	if err := json.Unmarshal(data, &redisUser); err != nil {
		return models.CachedUserDataOut{}, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return redisUser, nil
}
