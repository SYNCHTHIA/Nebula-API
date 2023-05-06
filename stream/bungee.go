package stream

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/nebulapb"
)

// PublishBungeeCommand - Send to proxy commands to BungeeCord
func PublishBungeeCommand(cmd string) error {
	c := pool.Get()
	defer c.Close()

	d := nebulapb.BungeeEntryStream{
		Type:    nebulapb.BungeeEntryStream_COMMAND,
		Command: cmd,
	}
	serialized, _ := json.Marshal(&d)
	_, err := c.Do("PUBLISH", "nebula.bungee.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Publish BungeeCommand")
		return err
	}
	return nil
}

// PublishBungeeEntry - PublishBungeeEntry Entry with Redis
func PublishBungeeEntry(data *nebulapb.BungeeEntry) error {
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
		logrus.WithError(err).Errorf("[Publish] Failed Publish BungeeEntry")
		return err
	}
	return nil
}
