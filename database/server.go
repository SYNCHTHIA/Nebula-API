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
	Fallback    bool          `bson:"fallback"`
	Status      PingResponse  `bson:"status"`
	//Status      StatusData    `json:"status" bson:"status"`
}

type PingResponse struct {
	Online      bool
	Version     VersionData       `json:"version"`
	Players     PlayersData       `json:"players"`
	Description map[string]string `json:"description"`
	Favicon     string            `json:"favicon"`
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
func GetAllServerEntry() ([]ServerData, error) {
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

// GetServerEntry - Get Individual Server Entry
func GetServerEntry(name string) (ServerData, error) {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("servers")

	server := ServerData{}
	err := coll.Find(bson.M{"name": name}).One(&server)
	if err != nil {
		return ServerData{}, err
	}

	return server, nil
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

// PushServerStatus - Push Server Status
func PushServerStatus(name string, response PingResponse) (string, int, error) {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("servers")

	v, err := coll.Find(bson.M{"name": name}).Count()

	// not found or nil
	if v == 0 || err != nil {
		return "", 0, err
	}

	updated := 0
	info, err := coll.Upsert(bson.M{"name": name}, bson.M{"$set": bson.M{"status": response}})
	if err == nil {
		updated = info.Updated
	}

	if updated > 0 {
		logrus.Debugf("[Fetcher] Updated : %s [%d]", name, updated)
	}

	return name, updated, err
}
