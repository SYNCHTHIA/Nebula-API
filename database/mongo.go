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

	m := &Mongo{
		client:   client,
		database: database,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		logrus.WithError(err).WithField("connection", mongoConStr).Fatalf("[MongoDB] Failed connect to MongoDB")
		return m
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.WithError(err).WithField("connection", mongoConStr).Errorf("[MongoDB] Failed ping to MongoDB")
		return m
	}

	logrus.Infof("[MongoDB] Connected to MongoDB")
	return m
}

func getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	return ctx, cancel
}
