package database

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type IPFilter struct {
	Id          string `gorm:"primaryKey;AutoIncrement;"`
	Action      FilterAction
	Address     string `gorm:"index;not null;"`
	Description string
}

type FilterAction int32

const (
	ALLOW FilterAction = iota
	DENY
)

func (s *Mysql) AddIPFilter(entry *IPFilter) error {
	r := s.client.First(&IPFilter{}, "address = ?", entry.Address)

	if r.RowsAffected != 0 {
		return errors.New("already exists")
	} else if r.Error != nil && r.RowsAffected != 0 {
		return r.Error
	}

	result := s.client.Create(entry)
	if result.Error != nil {
		logrus.WithError(result.Error).Errorf("[IPFW] Failed AddIPFilter")
		return result.Error
	}

	return nil
}

func (s *Mysql) RemoveIPFilter(address string) error {
	r := s.client.Delete(&IPFilter{}, "address = ?", address)
	return r.Error
}

func (s *Mysql) GetIPFilter(address string) (*IPFilter, error) {
	var entry *IPFilter
	r := s.client.Find(&entry, "address = ?", address)
	return entry, r.Error
}

// func (s *Mysql) ListIPFilter() ([]IPFilter, error) {
//     var entries []IPFilter
//
// }
