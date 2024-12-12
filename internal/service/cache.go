package service

import (
	"sync"
	"time"

	"github.com/cutlery47/email-service/internal/config"
	"github.com/cutlery47/email-service/internal/models"
)

// Кэш для хранения данных пользователей до момента подтверждения почты
type Cache interface {
	Put(user models.CachedUserData) error
	Get(mail string) (models.CachedUserData, error)
}

type MapCache struct {
	data map[string]models.CachedUserData
	mu   *sync.RWMutex

	conf config.Cache
}

func NewMapCache(conf config.Cache) *MapCache {
	cache := &MapCache{
		data: make(map[string]models.CachedUserData),
		mu:   &sync.RWMutex{},

		conf: conf,
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
