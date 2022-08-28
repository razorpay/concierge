package mutex

import (
	"time"

	"concierge/config"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var M *RedisMutexDriver

type RedisMutexDriver struct {
	New   *redsync.Redsync
	mutex *redsync.Mutex
}

// Pool ....
func Pool() *redsync.Redsync {
	log.Debug("Creating Connection")
	dbconfig := config.RedisDBConfig

	pool := redigo.NewPool(&redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle:     dbconfig.MaxIdle,
		IdleTimeout: 240 * time.Second,
		// max number of connections
		MaxActive: dbconfig.MaxActive,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				dbconfig.Host+":"+dbconfig.Port,
				redis.DialPassword(dbconfig.Password),
				redis.DialDatabase(dbconfig.Database))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	})
	return redsync.New(pool)
}

func (m *RedisMutexDriver) NewMutex(name string) {
	m.mutex = m.New.NewMutex(config.MutexPrefix + "_" + name)
}

func (m *RedisMutexDriver) Lock() error {
	return m.mutex.Lock()
}

func (m *RedisMutexDriver) Unlock() (bool, error) {
	return m.mutex.Unlock()
}
