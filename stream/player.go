package stream

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/nebulapb"
)

// PublishPlayerProfile - Stream for tablist
func PublishPlayerProfile(streamType nebulapb.PlayerPropertiesStream_Type, data *nebulapb.PlayerProfile) error {
	c := pool.Get()
	defer c.Close()

	if streamType != nebulapb.PlayerPropertiesStream_JOIN_SOLO && streamType != nebulapb.PlayerPropertiesStream_QUIT_SOLO {
		return errors.New("seriously?")
	}

	d := &nebulapb.PlayerPropertiesStream{
		Type: streamType,
		Solo: data,
	}

	serialized, _ := json.Marshal(&d)
	logrus.Debugln(d)

	_, err := c.Do("PUBLISH", "nebula.player.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Publish Player Profile")
		return err
	}
	return nil
}

// PublishAllPlayerProfile
func PublicAllPlayerProfile(streamType nebulapb.PlayerPropertiesStream_Type, data []*nebulapb.PlayerProfile) error {
	c := pool.Get()
	defer c.Close()

	if streamType != nebulapb.PlayerPropertiesStream_ADVERTISE_ALL {
		return errors.New("seriously?")
	}

	d := &nebulapb.PlayerPropertiesStream{
		Type: streamType,
		All:  data,
	}

	serialized, _ := json.Marshal(&d)
	logrus.Debugln(d)

	_, err := c.Do("PUBLISH", "nebula.player.global", string(serialized))
	if err != nil {
		logrus.WithError(err).Errorf("[Publish] Failed Publish Player Profile")
		return err
	}
	return nil
}
