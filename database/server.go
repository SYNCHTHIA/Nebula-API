package database

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// ServerData - Server List Data
type ServerData struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	DisplayName string        `json:"display_name" bson:"display_name"`
	Address     string        `json:"address" bson:"address"`
	Port        int32         `json:"port" bson:"port"`
	Motd        string        `json:"motd" bson:"motd"`
	Status      StatusData    `json:"status" bson:"status"`
}

// StatusData - Server Status
type StatusData struct {
	ServerStatus  string `bson:"server_status"`
	OnlinePlayers int32  `bson:"online_players"`
	MaxPlayers    int32  `bson:"max_players"`
}

// GetServerEntry - Get All Server Entries
func GetServerEntry() ([]ServerData, error) {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("servers")

	var servers []ServerData

	err := coll.Find(bson.M{}).All(&servers)
	if err != nil {
		logrus.WithError(err).Errorf("[Server] Failed Find ServerEntry: %s", err)
		return nil, err
	}

	return servers, nil
}

// AddServerEntry - Add Server Entry
func AddServerEntry(data ServerData) error {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("servers")

	if i, _ := coll.Find(bson.M{"name": data.Name}).Count(); i != 0 {
		return errors.New("already exists")
	}

	//err := coll.Insert(, data, bson.M{"status": &StatusData{}})

	err := coll.Insert(data)
	//_, err := coll.Upsert(bson.M{"name": data.Name}, bson.M{"status": &StatusData{}})
	if err != nil {
		logrus.WithError(err).Errorf("[Server] Failed AddServerEntry: %s", err)
	}

	return err
}

// RemoveServerEntry - RemoveServerEntry
func RemoveServerEntry(name string) error {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("servers")

	err := coll.Remove(bson.M{"name": name})
	if err != nil {
		logrus.WithError(err).Errorf("[Server] Failed RemoveServerEntry: %s", err)
	}

	return err
}
