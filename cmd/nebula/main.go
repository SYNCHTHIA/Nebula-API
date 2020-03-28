package main

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/database"
	"github.com/synchthia/nebula-api/logger"
	"github.com/synchthia/nebula-api/server"
	"github.com/synchthia/nebula-api/stream"
)

func startGRPC(port string, mongo *database.Mongo) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	return server.NewGRPCServer(mongo).Serve(lis)
}

func main() {
	// Init Logger
	logger.Init()

	// Init
	logrus.Printf("[NEBULA] Starting Nebula Server...")

	// Redis
	go func() {
		redisAddr := os.Getenv("REDIS_ADDRESS")
		if len(redisAddr) == 0 {
			redisAddr = "localhost:6379"
		}
		stream.NewRedisPool(redisAddr)
	}()

	// Connect to MongoDB
	mongoConStr := os.Getenv("MONGO_CONNECTION_STRING")
	if len(mongoConStr) == 0 {
		mongoConStr = "mongodb://localhost:27017"
	}

	mongoClient := database.NewMongoClient(mongoConStr, "nebula")

	// gRPC
	wait := make(chan struct{})
	go func() {
		defer close(wait)
		port := os.Getenv("GRPC_LISTEN_PORT")
		if len(port) == 0 {
			port = ":17200"
		}

		msg := logrus.WithField("listen", port)
		msg.Infof("[GRPC] Listening %s", port)

		if err := startGRPC(port, mongoClient); err != nil {
			logrus.Fatalf("[GRPC] gRPC Error: %s", err)
		}
	}()
	<-wait
}
