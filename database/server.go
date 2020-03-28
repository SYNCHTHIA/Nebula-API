package database

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/nebulapb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ServerData - Server List Data
type ServerData struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	Address     string             `json:"address" bson:"address"`
	Port        int32              `json:"port" bson:"port"`
	Motd        string             `json:"motd" bson:"motd"`
	Fallback    bool               `bson:"fallback"`
	Lockdown    Lockdown           `bson:"lockdown"`
	Status      PingResponse       `bson:"status"`
	//Status      StatusData    `json:"status" bson:"status"`
}

// Lockdown - Server lockdown entry
type Lockdown struct {
	Enabled     bool   `bson:"enabled"`
	Description string `bson:"description,omitempty"`
}

func LockdownFromProtobuf(pb *nebulapb.Lockdown) Lockdown {
	if pb == nil {
		return Lockdown{
			Enabled: false,
		}
	}

	return Lockdown{
		Enabled:     pb.Enabled,
		Description: pb.Description,
	}
}

func (l *Lockdown) ToProtobuf() *nebulapb.Lockdown {
	return &nebulapb.Lockdown{
		Enabled:     l.Enabled,
		Description: l.Description,
	}
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
func (m *Mongo) GetAllServerEntry() ([]ServerData, error) {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("servers")

	var servers []ServerData

	r, err := coll.Find(ctx, bson.M{})
	if err != nil {
		logrus.WithError(err).Errorf("[Server] Failed Find ServerEntry")
		return nil, err
	}

	if err := r.All(ctx, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

// GetServerEntry - Get Individual Server Entry
func (m *Mongo) GetServerEntry(name string) (ServerData, error) {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("servers")

	server := ServerData{}
	r := coll.FindOne(ctx, bson.M{"name": name})
	if r.Err() != nil {
		return ServerData{}, r.Err()
	}

	if err := r.Decode(&server); err != nil {
		return ServerData{}, err
	}

	return server, nil
}

// AddServerEntry - Add Server Entry
func (m *Mongo) AddServerEntry(data ServerData) error {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("servers")

	if cnt, err := coll.CountDocuments(ctx, bson.M{"name": data.Name}); cnt != 0 {
		return errors.New("already exists")
	} else if err != nil {
		return err
	}

	if _, err := coll.InsertOne(ctx, data); err != nil {
		logrus.WithError(err).Errorf("[Server] Failed AddServerEntry")
		return err
	}

	return nil
}

// RemoveServerEntry - RemoveServerEntry
func (m *Mongo) RemoveServerEntry(name string) error {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("servers")

	_, err := coll.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		logrus.WithError(err).Errorf("[Server] Failed RemoveServerEntry")
		return err
	}

	return nil
}

// PushServerStatus - Push Server Status
func (m *Mongo) PushServerStatus(name string, response PingResponse) (string, int, error) {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("servers")

	if cnt, err := coll.CountDocuments(ctx, bson.M{"name": name}); cnt == 0 || err != nil {
		return "", 0, err
	}

	updated := 0
	r, err := coll.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"status": response}}, options.Update().SetUpsert(true))
	if err == nil && r.ModifiedCount >= 1 {
		updated = int(r.ModifiedCount)
		logrus.Debugf("[Fetcher] Updated : %s [%d]", name, updated)
	}

	return name, updated, err
}

// SetLockdown - Set server Lockdown
func (m *Mongo) SetLockdown(name string, enabled bool, description string) error {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("servers")

	lockdown := &Lockdown{
		Enabled:     enabled,
		Description: description,
	}

	_, err := coll.UpdateOne(ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"lockdown": lockdown}}, options.Update().SetUpsert(true))

	return err
}
