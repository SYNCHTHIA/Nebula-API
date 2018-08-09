package database

import (
	"github.com/globalsign/mgo/bson"
)

// BungeeData - BungeeCord configuration Data
type BungeeData struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Motd    string        `bson:"motd"`
	Favicon string        `bson:"favicon"`
}

// GetBungeeEntry - Get Bungee Entry
func GetBungeeEntry() (BungeeData, error) {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("bungee")

	bungee := BungeeData{}
	err := coll.Find(bson.M{}).One(&bungee)
	if err != nil {
		return BungeeData{}, err
	}

	return bungee, nil
}

// SetMotd - Set Motd
func SetMotd(motd string) error {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("bungee")

	_, err := coll.Upsert(bson.M{}, bson.M{"$set": bson.M{"motd": motd}})
	return err
}

// SetFavicon - Set Favicon
func SetFavicon(favicon string) error {
	session := session.Copy()
	defer session.Close()
	coll := session.DB("nebula").C("bungee")

	_, err := coll.Upsert(bson.M{}, bson.M{"$set": bson.M{"favicon": favicon}})
	return err
}
