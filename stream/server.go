package stream

import (
	"encoding/json"

	"github.com/synchthia/nebula-api/nebulapb"

	"github.com/sirupsen/logrus"
)

// PublishServer - Publish with Redis
func PublishServer(data *nebulapb.ServerEntry) {
	c := pool.Get()
	defer c.Close()

	d := &nebulapb.ServerEntryStream{
		Type:  nebulapb.ServerEntryStream_SYNC,
		Entry: data,
	}
	serialized, _ := json.Marshal(&d)
	//logrus.Debugln(d)

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
	//logrus.Debugln(data)
	_, err := c.Do("PUBLISH", "nebula.servers.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Remove Server")
	}
}
