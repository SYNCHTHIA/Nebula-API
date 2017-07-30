package server

import (
	"sync"

	"golang.org/x/net/context"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Nebula-API/database"
	pb "gitlab.com/Startail/Nebula-API/nebulapb"
	"google.golang.org/grpc"
)

type Server interface {
	Ping()
	GetServerEntry()
	AddServerEntry()
	RemoveServerEntry()
	PushStatus()
	FetchStatus()
}

type grpcServer struct {
	server           Server
	mu               sync.RWMutex
	entryStreamChans map[chan pb.EntryStreamResponse]struct{}
}

func NewServer() *grpcServer {
	return &grpcServer{
		entryStreamChans: make(map[chan pb.EntryStreamResponse]struct{}),
	}
}

func NewGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterNebulaServer(server, NewServer())
	return server
}

func (s *grpcServer) Ping(ctx context.Context, e *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *grpcServer) GetServerEntry(ctx context.Context, e *pb.GetServerEntryRequest) (*pb.GetServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var rpcServerEntry []*pb.ServerEntry

	db, err := database.GetServerEntry()
	if err != nil {
		logrus.WithError(err).Errorf("[gRPC] Error @ GetServerEntry: %s", err)
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

	dbEntry := s.ServerEntry_PBtoDB(*e.Entry)
	err := database.AddServerEntry(dbEntry)

	// Stream Casting
	for c := range s.entryStreamChans {
		c <- pb.EntryStreamResponse{Type: pb.StreamType_SYNC, Target: "BungeeCord"}
	}
	return &pb.AddServerEntryResponse{}, err
}

func (s *grpcServer) RemoveServerEntry(ctx context.Context, e *pb.RemoveServerEntryRequest) (*pb.RemoveServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := database.RemoveServerEntry(e.Name)

	// Stream Casting
	for c := range s.entryStreamChans {
		c <- pb.EntryStreamResponse{Type: pb.StreamType_SYNC, Target: "BungeeCord"}
	}
	return &pb.RemoveServerEntryResponse{}, err
}

func (s *grpcServer) PushStatus(ctx context.Context, e *pb.PushStatusRequest) (*pb.PushStatusResponse, error) {
	return &pb.PushStatusResponse{}, nil
}

func (s *grpcServer) FetchStatus(ctx context.Context, e *pb.FetchStatusRequest) (*pb.FetchStatusResponse, error) {
	return &pb.FetchStatusResponse{}, nil
}

func (s *grpcServer) ServerEntry_DBtoPB(dbEntry database.ServerData) *pb.ServerEntry {
	return &pb.ServerEntry{
		Name:        dbEntry.Name,
		DisplayName: dbEntry.DisplayName,
		Address:     dbEntry.Address,
		Port:        dbEntry.Port,
		Motd:        dbEntry.Motd,
	}
}

func (s *grpcServer) ServerEntry_PBtoDB(pbEntry pb.ServerEntry) database.ServerData {
	return database.ServerData{
		Name:        pbEntry.Name,
		DisplayName: pbEntry.DisplayName,
		Address:     pbEntry.Address,
		Port:        pbEntry.Port,
		Motd:        pbEntry.Motd,
	}
}
