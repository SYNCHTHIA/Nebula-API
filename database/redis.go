package database

import (
	"encoding/json"
	"time"

	"gitlab.com/Startail/Nebula-API/nebulapb"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

var pool *redis.Pool
var conn redis.Conn

// NewRedisPool - redis Connection Pooling
func NewRedisPool(server string) {
	logrus.WithFields(logrus.Fields{
		"server": server,
	}).Infof("[Redis] Creating Pool...")

	pool = &redis.Pool{
		MaxIdle:   12,
		MaxActive: 0,
		//IdleTimeout: 240 * time.Second,
		Wait: true,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)

			if err != nil {
				logrus.WithError(err).Errorf("[Redis] Error occurred in Connecting: %s", server)
				return nil, err
			}

			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				logrus.WithError(err).Errorf("[Redis] Error occurred in Redis Pool: %s", server)
			}
			return err
		},
	}
}

// PublishServer - Publish with Redis
func PublishServer(data *nebulapb.ServerEntry) {
	c := pool.Get()
	defer c.Close()

	d := &nebulapb.ServerEntryStream{
		Type:  nebulapb.ServerEntryStream_SYNC,
		Entry: data,
	}
	serialized, _ := json.Marshal(&d)
	logrus.Debugln(d)

	_, err := c.Do("PUBLISH", "nebula.servers.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Publish Server")
	}
}

// PublishRemoveServer - Remove Server send Redis
func PublishRemoveServer(data *nebulapb.ServerEntry) {
	c := pool.Get()
	defer c.Close()

	d := &nebulapb.ServerEntryStream{
		Type:  nebulapb.ServerEntryStream_REMOVE,
		Entry: data,
	}
	serialized, _ := json.Marshal(&d)
	logrus.Debugln(data)
	_, err := c.Do("PUBLISH", "nebula.servers.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Remove Server")
	}
}
