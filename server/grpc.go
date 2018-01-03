package server

import (
	"sync"

	"golang.org/x/net/context"

	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Nebula-API/database"
	pb "gitlab.com/Startail/Nebula-API/nebulapb"
	"google.golang.org/grpc"
)

type Server interface {
	GetAllServerEntry()
	AddServerEntry()
	RemoveServerEntry()
	//FetchStatus()
}

type grpcServer struct {
	server Server
	mu     sync.RWMutex
}

func NewServer() *grpcServer {
	return &grpcServer{}
}

func NewGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	newServer := NewServer()
	pb.RegisterNebulaServer(server, newServer)

	// Pinging
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				newServer.pinging()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return server
}

func (s *grpcServer) GetServerEntry(ctx context.Context, e *pb.GetServerEntryRequest) (*pb.GetServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var rpcServerEntry []*pb.ServerEntry

	db, err := database.GetAllServerEntry()
	if err != nil {
		logrus.WithError(err).Errorf("[gRPC] Error @ GetAllServerEntry: %s", err)
		return nil, err
	}
	for _, ent := range db {
		pbEntry := s.ServerEntry_DBtoPB(ent)
		rpcServerEntry = append(rpcServerEntry, pbEntry)
	}
	return &pb.GetServerEntryResponse{Entry: rpcServerEntry}, err
}

func (s *grpcServer) AddServerEntry(ctx context.Context, e *pb.AddServerEntryRequest) (*pb.AddServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dbEntry := s.ServerEntry_PBtoDB(e.Entry)
	err := database.AddServerEntry(dbEntry)

	database.PublishServer(e.Entry)

	return &pb.AddServerEntryResponse{}, err
}

func (s *grpcServer) RemoveServerEntry(ctx context.Context, e *pb.RemoveServerEntryRequest) (*pb.RemoveServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := database.RemoveServerEntry(e.Name)

	if err == nil {
		database.PublishRemoveServer(&pb.ServerEntry{Name: e.Name})
	}
	return &pb.RemoveServerEntryResponse{}, err
}

func (s *grpcServer) Status_DBtoPB(dbEntry database.PingResponse) *pb.ServerStatus {
	return &pb.ServerStatus{
		Online: dbEntry.Online,
		Version: &pb.ServerStatus_Version{
			Name:     dbEntry.Version.Name,
			Protocol: int32(dbEntry.Version.Protocol),
		},
		Players: &pb.ServerStatus_Players{
			Max:    int32(dbEntry.Players.Max),
			Online: int32(dbEntry.Players.Online),
		},
		Description: dbEntry.Description["text"],
		Favicon:     dbEntry.Favicon,
	}
}

func (s *grpcServer) ServerEntry_DBtoPB(dbEntry database.ServerData) *pb.ServerEntry {
	return &pb.ServerEntry{
		Name:        dbEntry.Name,
		DisplayName: dbEntry.DisplayName,
		Address:     dbEntry.Address,
		Port:        dbEntry.Port,
		Motd:        dbEntry.Motd,
		Fallback:    dbEntry.Fallback,
		Status:      s.Status_DBtoPB(dbEntry.Status),
	}
}

func (s *grpcServer) ServerEntry_PBtoDB(pbEntry *pb.ServerEntry) database.ServerData {
	return database.ServerData{
		Name:        pbEntry.Name,
		DisplayName: pbEntry.DisplayName,
		Address:     pbEntry.Address,
		Port:        pbEntry.Port,
		Motd:        pbEntry.Motd,
		Fallback:    pbEntry.Fallback,
	}
}
