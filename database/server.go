package database

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// Servers - Server definition
type Servers struct {
	Id          int32  `gorm:"primaryKey;AutoIncrement;"`
	Name        string `gorm:"index;not null;"`
	DisplayName string
	Address     string
	Port        int32
	Motd        string
	Fallback    bool
	//Lockdown    *Lockdown `gorm:"references:Lockdown"`
	Lockdown string `gorm:"type:json;"`
	Status   string `gorm:"type:json;"`
}

// Lockdown - Server lockdown entry
type Lockdown struct {
	Id          uint
	ServersId   int32
	Enabled     bool
	Description string
}

type PingResponse struct {
	Online      bool
	Version     VersionData `json:"version"`
	Players     PlayersData `json:"players"`
	Description interface{} `json:"description"`
	Favicon     string      `json:"favicon"`
}

type VersionData struct {
	Name     string
	Protocol int
}

type PlayersData struct {
	Max    int32
	Online int32
}

// StatusData - Server Status
/*type StatusData struct {
	ServerStatus  string `bson:"server_status"`
	OnlinePlayers int32  `bson:"online_players"`
	MaxPlayers    int32  `bson:"max_players"`
}*/

// GetAllServerEntry - Get All Server Entries
func (s *Mysql) GetAllServerEntry() ([]Servers, error) {
	var servers []Servers
	r := s.client.Find(&servers)
	if r.Error != nil {
		logrus.WithError(r.Error).Errorf("[Server] Failed Find ServerEntry")
		return nil, r.Error
	}

	return servers, nil
}

// GetServerEntry - Get Individual Server Entry
func (s *Mysql) GetServerEntry(name string) (Servers, error) {
	server := Servers{}
	r := s.client.Find(&server, "name = ?", name)
	if r.Error != nil {
		return Servers{}, r.Error
	}

	return server, nil
}

// AddServerEntry - Add Server Entry
func (s *Mysql) AddServerEntry(data Servers) error {
	r := s.client.First(&Servers{}, "name = ?", data.Name)

	if r.RowsAffected != 0 {
		return errors.New("already exists")
	} else if r.Error != nil && r.RowsAffected != 0 {
		return r.Error
	}

	result := s.client.Create(&Servers{
		Name:        data.Name,
		DisplayName: data.DisplayName,
		Address:     data.Address,
		Port:        data.Port,
		Motd:        data.Motd,
		Fallback:    data.Fallback,
		Lockdown:    data.Lockdown,
		Status:      "{}",
	})

	if result.Error != nil {
		logrus.WithError(result.Error).Errorf("[Server] Failed AddServerEntry")
		return result.Error
	}

	return nil
}

// RemoveServerEntry - RemoveServerEntry
func (s *Mysql) RemoveServerEntry(name string) error {
	r := s.client.Delete(&Servers{}, "name = ?", name)
	if r.Error != nil {
		logrus.WithError(r.Error).Errorf("[Server] Failed RemoveServerEntry")
		return r.Error
	}

	return nil
}

// PushServerStatus - Push Server Status
func (s *Mysql) PushServerStatus(name, response string) (string, int, error) {
	r := s.client.Model(&Servers{}).Where("name = ?", name).Update("status", response)

	if r.Error != nil {
		return "", 0, r.Error
	}

	return name, 0, r.Error
}

// SetLockdown - Set server Lockdown
func (s *Mysql) SetLockdown(name string, enabled bool, description string) error {
	r := s.client.Model(&Servers{}).Where("name = ?", name).Update("lockdown", Lockdown{
		Enabled:     enabled,
		Description: description,
	})

	return r.Error
}
