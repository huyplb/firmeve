package cache

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/kataras/iris/core/errors"
	"strings"
	"sync"
	"time"

	"github.com/firmeve/firmeve/cache/redis"
	goRedis "github.com/go-redis/redis"
)

type Cache struct {
	config       contract.Configuration
	current      contract.CacheSerializable
	repositories map[string]contract.CacheSerializable
}

var (
	mutex             sync.Mutex
	ErrDriverNotFound = errors.New(`driver not found`)
)

// Create a cache manager
func New(config contract.Configuration) contract.Cache {
	cache := &Cache{
		config:       config,
		repositories: make(map[string]contract.CacheSerializable, 0),
	}
	cache.registerDefaultDriver(cache.config.GetString(`default`))
	cache.current = cache.Driver(cache.config.GetString(`default`))

	return cache
}

// Get the cache driver of the finger
func (c *Cache) Driver(driver string) contract.CacheSerializable {
	var current contract.CacheSerializable
	var ok bool

	mutex.Lock()
	defer mutex.Unlock()

	if current, ok = c.repositories[driver]; ok {
		return current
	}

	panic(ErrDriverNotFound)
}

// Register driver
func (c *Cache) Register(driver string, store contract.CacheStore) {
	c.repositories[driver] = NewRepository(store)
}

// Create a redis cache driver
func (c *Cache) createRedisDriver() contract.CacheStore {
	var (
		host   = c.config.GetString(`repositories.redis.host`)
		port   = c.config.GetString(`repositories.redis.port`)
		db     = c.config.GetInt(`repositories.redis.db`)
		prefix = c.config.GetString(`prefix`)
	)

	addr := []string{host, `:`, port}

	return redis.New(goRedis.NewClient(&goRedis.Options{
		Addr: strings.Join(addr, ``),
		DB:   db,
	}), prefix)
}

// register a exists default driver
func (c *Cache) registerDefaultDriver(driver string) {
	var store contract.CacheStore
	switch driver {
	case `redis`:
		store = c.createRedisDriver()
	default:
		panic(ErrDriverNotFound)
	}

	c.Register(driver, store)
}

func (c *Cache) Store() contract.CacheStore {
	return c.current.Store()
}

func (c *Cache) Get(key string) (interface{}, error) {
	return c.current.Store().Get(key)
}

func (c *Cache) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	return c.current.GetDefault(key, defaultValue)
}

func (c *Cache) Pull(key string) (interface{}, error) {
	return c.current.Pull(key)
}

func (c *Cache) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	return c.current.PullDefault(key, defaultValue)
}

func (c *Cache) Add(key string, value interface{}, expire time.Time) error {
	return c.current.Store().Add(key, value, expire)
}

func (c *Cache) Put(key string, value interface{}, expire time.Time) error {
	return c.current.Store().Put(key, value, expire)
}

func (c *Cache) Forever(key string, value interface{}) error {
	return c.current.Store().Forever(key, value)
}

func (c *Cache) Forget(key string) error {
	return c.current.Store().Forget(key)
}

func (c *Cache) Increment(key string, steps ...int64) error {
	return c.current.Store().Increment(key, steps...)
}

func (c *Cache) Decrement(key string, steps ...int64) error {
	return c.current.Store().Decrement(key, steps...)
}

func (c *Cache) Has(key string) bool {
	return c.current.Store().Has(key)
}

func (c *Cache) Flush() error {
	return c.current.Store().Flush()
}

func (c *Cache) GetDecode(key string, to interface{}) (interface{}, error) {
	return c.current.GetDecode(key, to)
}

func (c *Cache) AddEncode(key string, value interface{}, expire time.Time) error {
	return c.current.AddEncode(key, value, expire)
}

func (c *Cache) ForeverEncode(key string, value interface{}) error {
	return c.current.ForeverEncode(key, value)
}

func (c *Cache) PutEncode(key string, value interface{}, expire time.Time) error {
	return c.current.PutEncode(key, value, expire)
}
