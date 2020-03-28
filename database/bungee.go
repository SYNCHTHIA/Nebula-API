package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BungeeData - BungeeCord configuration Data
type BungeeData struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Motd    string             `bson:"motd"`
	Favicon string             `bson:"favicon"`
}

// GetBungeeEntry - Get Bungee Entry
func (m *Mongo) GetBungeeEntry() (BungeeData, error) {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("bungee")

	bungee := BungeeData{}
	r := coll.FindOne(ctx, bson.M{})
	if r.Err() != nil {
		return BungeeData{}, r.Err()
	}

	if err := r.Decode(&bungee); err != nil {
		return BungeeData{}, err
	}

	return bungee, nil
}

// SetMotd - Set Motd
func (m *Mongo) SetMotd(motd string) error {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("bungee")

	r := coll.FindOneAndUpdate(ctx, bson.M{}, bson.M{"$set": bson.M{"motd": motd}}, options.FindOneAndUpdate().SetUpsert(true))
	if r.Err() != nil {
		return r.Err()
	}

	return nil
}

// SetFavicon - Set Favicon
func (m *Mongo) SetFavicon(favicon string) error {
	ctx, cancel := getContext()
	defer cancel()
	coll := m.client.Database(m.database).Collection("bungee")

	r := coll.FindOneAndUpdate(ctx, bson.M{}, bson.M{"$set": bson.M{"favicon": favicon}}, options.FindOneAndUpdate().SetUpsert(true))
	if r.Err() != nil {
		return r.Err()
	}

	return nil
}
