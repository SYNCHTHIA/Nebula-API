package server

import (
	"encoding/json"
	"errors"
	"sync"

	"golang.org/x/net/context"

	"time"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/database"
	pb "github.com/synchthia/nebula-api/nebulapb"
	"github.com/synchthia/nebula-api/service"
	"github.com/synchthia/nebula-api/stream"
	"google.golang.org/grpc"
)

type Services struct {
	MySQL    *database.Mysql
	IPFilter *service.IPFilter
}

type Server interface {
	GetAllServerEntry()
	AddServerEntry()
	RemoveServerEntry()
	GetBungeeEntry()
	SetLockdown()
	//FetchStatus()
}

type grpcServer struct {
	server Server
	mu     sync.RWMutex
	svc    *Services
}

func NewServer(svc *Services) *grpcServer {
	return &grpcServer{
		svc: svc,
	}
}

func NewGRPCServer(svc *Services) *grpc.Server {
	server := grpc.NewServer()
	newServer := NewServer(svc)
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

	db, err := s.svc.MySQL.GetAllServerEntry()
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
	err := s.svc.MySQL.AddServerEntry(dbEntry)

	stream.PublishServer(e.Entry)

	return &pb.AddServerEntryResponse{}, err
}

func (s *grpcServer) RemoveServerEntry(ctx context.Context, e *pb.RemoveServerEntryRequest) (*pb.RemoveServerEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.svc.MySQL.RemoveServerEntry(e.Name)

	if err == nil {
		stream.PublishRemoveServer(&pb.ServerEntry{Name: e.Name})
	}
	return &pb.RemoveServerEntryResponse{}, err
}

func (s *grpcServer) GetBungeeEntry(ctx context.Context, e *pb.GetBungeeEntryRequest) (*pb.GetBungeeEntryResponse, error) {
	entry, err := s.svc.MySQL.GetBungeeEntry()
	return &pb.GetBungeeEntryResponse{Entry: s.BungeeEntry_DBtoPB(entry)}, err
}

func (s *grpcServer) SetMotd(ctx context.Context, e *pb.SetMotdRequest) (*pb.SetMotdResponse, error) {
	err := s.svc.MySQL.SetMotd(e.Motd)
	entry, err := s.svc.MySQL.GetBungeeEntry()
	stream.PublishBungee(s.BungeeEntry_DBtoPB(entry))
	return &pb.SetMotdResponse{}, err
}

func (s *grpcServer) SetFavicon(ctx context.Context, e *pb.SetFaviconRequest) (*pb.SetFaviconResponse, error) {
	err := s.svc.MySQL.SetFavicon(e.Favicon)
	entry, err := s.svc.MySQL.GetBungeeEntry()
	stream.PublishBungee(s.BungeeEntry_DBtoPB(entry))
	return &pb.SetFaviconResponse{}, err
}

func (s *grpcServer) SetLockdown(ctx context.Context, e *pb.SetLockdownRequest) (*pb.SetLockdownResponse, error) {
	if e.Lockdown.Enabled && e.Lockdown.Description == "" {
		e.Lockdown.Description = "&cThis server currently not available"
	}

	if err := s.svc.MySQL.SetLockdown(e.Name, e.Lockdown.Enabled, e.Lockdown.Description); err != nil {
		return &pb.SetLockdownResponse{}, err
	}

	entry, err := s.svc.MySQL.GetServerEntry(e.Name)
	if err == nil {
		stream.PublishServer(s.ServerEntry_DBtoPB(entry))
	}

	return &pb.SetLockdownResponse{Entry: s.ServerEntry_DBtoPB(entry)}, err
}

func (s *grpcServer) IPLookup(ctx context.Context, e *pb.IPLookupRequest) (*pb.IPLookupResponse, error) {
	if s.svc.IPFilter != nil {
		res, err := s.svc.IPFilter.Check(e.IpAddress)
		if err != nil {
			return &pb.IPLookupResponse{}, err
		}

		return &pb.IPLookupResponse{
			Result: &pb.IPLookupResult{
				IpAddress:    res.IPAddress,
				Isp:          res.ISP,
				IsSuspicious: res.IsSuspicious,
				Reason:       res.Reason,
			},
		}, nil
	} else {
		return &pb.IPLookupResponse{}, errors.New("iplookup not enabled")
	}
}

func (s *grpcServer) BungeeEntry_DBtoPB(dbEntry database.Bungee) *pb.BungeeEntry {
	return &pb.BungeeEntry{
		Motd:    dbEntry.Motd,
		Favicon: dbEntry.Favicon,
	}
}

func (s *grpcServer) BungeeEntry_PBtoDB(pbEntry *pb.BungeeEntry) database.Bungee {
	return database.Bungee{
		Motd:    pbEntry.Motd,
		Favicon: pbEntry.Favicon,
	}
}

func (s *grpcServer) Lockdown_DBtoPB(dbEntry database.Lockdown) *pb.Lockdown {
	return &pb.Lockdown{
		Enabled:     dbEntry.Enabled,
		Description: dbEntry.Description,
	}
}

func (s *grpcServer) Lockdown_PBtoDB(pbEntry *pb.Lockdown) *database.Lockdown {
	if pbEntry == nil {
		return &database.Lockdown{
			Enabled: false,
		}
	}

	return &database.Lockdown{
		Enabled:     pbEntry.Enabled,
		Description: pbEntry.Description,
	}
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

func (s *grpcServer) ServerEntry_DBtoPB(dbEntry database.Servers) *pb.ServerEntry {
	lockdown := database.Lockdown{}
	json.Unmarshal([]byte(dbEntry.Lockdown), &lockdown)

	status := database.PingResponse{}
	json.Unmarshal([]byte(dbEntry.Status), &status)

	return &pb.ServerEntry{
		Name:        dbEntry.Name,
		DisplayName: dbEntry.DisplayName,
		Address:     dbEntry.Address,
		Port:        dbEntry.Port,
		Motd:        dbEntry.Motd,
		Fallback:    dbEntry.Fallback,
		Lockdown:    s.Lockdown_DBtoPB(lockdown),
		Status:      s.Status_DBtoPB(status),
	}
}

func (s *grpcServer) ServerEntry_PBtoDB(pbEntry *pb.ServerEntry) database.Servers {
	r, _ := json.Marshal(s.Lockdown_PBtoDB(pbEntry.Lockdown))

	return database.Servers{
		Name:        pbEntry.Name,
		DisplayName: pbEntry.DisplayName,
		Address:     pbEntry.Address,
		Port:        pbEntry.Port,
		Motd:        pbEntry.Motd,
		Fallback:    pbEntry.Fallback,
		Lockdown:    string(r),
	}
}
