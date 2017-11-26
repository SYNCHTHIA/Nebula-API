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
	EntryDispatch()
	Ping()
	GetAllServerEntry()
	AddServerEntry()
	RemoveServerEntry()
	//FetchStatus()
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

func (s *grpcServer) Ping(ctx context.Context, e *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
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

	for c := range s.entryStreamChans {
		c <- pb.EntryStreamResponse{Type: pb.StreamType_SYNC, Entry: e.Entry}
	}
	return &pb.AddServerEntryResponse{}, err
}

func (s *grpcServer) RemoveServerEntry(ctx context.Context, e *pb.RemoveServerEntryRequest) (*pb.RemoveServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := database.RemoveServerEntry(e.Name)

	if err == nil {
		for c := range s.entryStreamChans {
			// Announce only Name for Delete local entry
			c <- pb.EntryStreamResponse{Type: pb.StreamType_REMOVE, Entry: &pb.ServerEntry{Name: e.Name}}
		}

	}
	return &pb.RemoveServerEntryResponse{}, err
}

/*
func (s *grpcServer) PushStatus(ctx context.Context, e *pb.PushStatusRequest) (*pb.PushStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//todo: get is updated
	updated, err := database.PushServerStatus(e.Name, e.Status)
	if updated > 0 {
		//todo: if updated, send to spigot?
		// -> get server entry from name
		logrus.Printf("[PushStatus] Changed Something!")
	}

	return &pb.PushStatusResponse{}, err
} */

/*func (s *grpcServer) FetchStatus(ctx context.Context, e *pb.FetchStatusRequest) (*pb.FetchStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	logrus.Printf("Incoming Request")

	var rpcServerEntry []*pb.ServerEntry

	db, err := database.GetAllServerEntry()
	if err != nil {
		logrus.WithError(err).Errorf("[gRPC] Error @ FetchStatus: %s", err)
		return nil, err
	}
	for _, ent := range db {
		logrus.Printf("[***] %s", ent.Name)
		pbEntry := s.ServerEntry_DBtoPB(ent)
		rpcServerEntry = append(rpcServerEntry, pbEntry)
	}
	return &pb.FetchStatusResponse{Entry: rpcServerEntry}, nil
}*/

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
	}
}
