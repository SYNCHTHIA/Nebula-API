package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Servers struct {
	Id          int32  `gorm:"primaryKey;AutoIncrement;"`
	Name        string `gorm:"index;not null;"`
	DisplayName string
	Address     string
	Port        int32
	Motd        string
	Fallback    bool
	Lockdown    string
	Status      string
}

type Bungee struct {
	Id      int32 `gorm:"primaryKey;AutoIncrement;"`
	Motd    string
	Favicon string
}

type Mysql struct {
	client   *gorm.DB
	database string
}

func NewMysqlClient(mysqlConStr, database string) *Mysql {
	logrus.WithField("connection", mysqlConStr).Infof("[MySQL] Connecting to MySQL...")

	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       mysqlConStr,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("[MySQL] Failed to connect to MySQL: %s", err)
		return nil
	}

	m := &Mysql{
		client:   client,
		database: database,
	}

	if err := m.client.AutoMigrate(&Servers{}); err != nil {
		logrus.Fatalf("[MySQL] Failed to migrate: %s", err)
		return m
	}

	if err := m.client.AutoMigrate(&Bungee{}); err != nil {
		logrus.Fatalf("[MySQL] Failed to migrate: %s", err)
		return m
	}

	logrus.Infof("[MySQL] Connected to MySQL")

	return m
}
