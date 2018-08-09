package stream

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Nebula-API/nebulapb"
)

// PublishBungee - PublishBungee Entry with Redis
func PublishBungee(data *nebulapb.BungeeEntry) {
	c := pool.Get()
	defer c.Close()

	d := &nebulapb.BungeeEntryStream{
		Type:  nebulapb.BungeeEntryStream_SYNC,
		Entry: data,
	}
	serialized, _ := json.Marshal(&d)
	logrus.Debugln(d)

	_, err := c.Do("PUBLISH", "nebula.bungee.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Publish Bungee")
	}
}
