package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	client   *mongo.Client
	database string
}

// NewMongoClient - New MongoDB Client
func NewMongoClient(mongoConStr, database string) *Mongo {
	logrus.WithField("connection", mongoConStr).Infof("[MongoDB] Connecting to MongoDB...")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConStr))
	if err != nil {
		logrus.WithError(err).WithField("connection", mongoConStr).Fatalf("[MongoDB] Failed ensure Client of MongoDB")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logrus.WithError(err).WithField("connection", mongoConStr).Fatalf("[MongoDB] Failed connect to MongoDB")
		return nil
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logrus.WithError(err).WithField("connection", mongoConStr).Fatalf("[MongoDB] Failed ping to MongoDB")
		return nil
	}

	logrus.Infof("[MongoDB] Connected to MongoDB")

	return &Mongo{
		client:   client,
		database: database,
	}
}

func getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}
