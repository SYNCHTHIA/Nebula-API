package database

import (
	"encoding/json"
	"flag"
	"strings"
	"time"

	"gitlab.com/Startail/Nebula-API/nebulapb"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

var redisConn redis.Conn

// NewRedisSession - Connect to Redis
func NewRedisSession(address string) {
	// Set Flag
	pAddress := strings.Split(address, ":")
	host := flag.String("hostname", pAddress[0], "Set Hostname")
	port := flag.String("port", pAddress[1], "Set Port")
	flag.Parse()

	for {
		logrus.WithFields(logrus.Fields{
			"servers": address,
		}).Infof("[Redis] Connecting...")

		// Connecting
		c, err := redis.Dial("tcp", *host+":"+*port)
		if err != nil {
			logrus.WithError(err).Errorf("[Redis] Error occurred while connecting to %s", address)
			time.Sleep(5 * time.Second)
			continue
		}

		redisConn = c
		logrus.Printf("[Redis] Connected!")
		for {
			if c.Err() != nil {
				logrus.WithError(err).Errorf("[Redis] Error occurred in Connected Session: %s", address)
				c.Close()
				time.Sleep(5 * time.Second)
				break
			}
		}
	}
}

// DisconnectRedis - Disconnect from Redis
func DisconnectRedis() {
	redisConn.Close()
}

// PublishServer - Publish with Redis
func PublishServer(data *nebulapb.ServerEntry) {
	d := &nebulapb.ServerEntryStream{
		Type:  nebulapb.ServerEntryStream_SYNC,
		Entry: data,
	}
	serialized, _ := json.Marshal(&d)
	logrus.Debugln(d)

	_, err := redisConn.Do("PUBLISH", "nebula.servers.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Publish Server")
	}
}

// PublishRemoveServer - Remove Server send Redis
func PublishRemoveServer(data *nebulapb.ServerEntry) {
	d := &nebulapb.ServerEntryStream{
		Type:  nebulapb.ServerEntryStream_REMOVE,
		Entry: data,
	}
	serialized, _ := json.Marshal(&d)
	logrus.Debugln(data)
	_, err := redisConn.Do("PUBLISH", "nebula.servers.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Remove Server")
	}
}
