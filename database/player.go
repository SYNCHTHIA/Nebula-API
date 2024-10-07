package database

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/nebulapb"
	"gorm.io/gorm/clause"
)

type Players struct {
	ID            uint   `gorm:"primary_key;AutoIncrement;"`
	UUID          string `gorm:"index;unique;"`
	Name          string `gorm:"index;not null;"`
	CurrentServer string
	Latency       int64
	RawProperties string `gorm:"type:text"`
}

type UpdateOption struct {
	IsQuit bool
}

func (p *Players) ToProtobuf() *nebulapb.PlayerProfile {
	return &nebulapb.PlayerProfile{
		PlayerUUID:    p.UUID,
		PlayerName:    p.Name,
		PlayerLatency: int64(p.Latency),
		CurrentServer: p.CurrentServer,
		Properties: func() []*nebulapb.PlayerProperty {
			properties := []*nebulapb.PlayerProperty{}
			if err := json.Unmarshal([]byte(p.RawProperties), &properties); err != nil {
				return nil
			}
			return properties
		}(),
	}
}

func PlayersFromProtobuf(p *nebulapb.PlayerProfile) *Players {
	return &Players{
		UUID:          p.PlayerUUID,
		Name:          p.PlayerName,
		Latency:       int64(p.PlayerLatency),
		CurrentServer: p.CurrentServer,
		RawProperties: func() string {
			if b, err := json.Marshal(p.Properties); err != nil {
				return "[]"
			} else {
				return string(b)
			}
		}(),
	}
}

func (s *Mysql) GetAllPlayers() ([]Players, error) {
	var players []Players
	r := s.client.Where("current_server != ?", "").Find(&players)
	if r.Error != nil {
		logrus.WithError(r.Error).Errorf("[Player] Failed Find Player")
		return nil, r.Error
	}

	return players, nil
}

func (s *Mysql) UpdateAllPlayers(players []Players) error {
	r := s.client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		UpdateAll: true,
	}).Create(players)

	return r.Error
}

func (s *Mysql) SyncPlayer(newPlayer *Players, opts *UpdateOption) error {
	var player Players
	findRes := s.client.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&Players{}).First(&player, "uuid = ?", newPlayer.UUID)
	if findRes.Error != nil {
		logrus.WithError(findRes.Error).Errorf("[Player] SyncPlayer: Failed update player data (%s)", newPlayer.UUID)
		return findRes.Error
	}

	if opts.IsQuit {
		newPlayer.CurrentServer = ""
	} else {
		player.CurrentServer = newPlayer.CurrentServer
	}

	r := s.client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		UpdateAll: true,
	}).Create(newPlayer)

	return r.Error
}
