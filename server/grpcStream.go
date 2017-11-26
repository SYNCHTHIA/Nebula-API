package server

import (
	"golang.org/x/net/context"

	"strings"

	"github.com/sirupsen/logrus"
	pb "gitlab.com/Startail/Nebula-API/nebulapb"
)

func (s *grpcServer) QuitEntryStream(ctx context.Context, e *pb.QuitEntryStreamRequest) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for c := range s.entryStreamChans {
		c <- pb.EntryStreamResponse{Target: e.Name, Type: pb.StreamType_QUIT}
	}

	return &pb.Empty{}, nil
}

func (s *grpcServer) EntryStream(r *pb.StreamRequest, es pb.Nebula_EntryStreamServer) error {
	ech := make(chan pb.EntryStreamResponse)
	s.mu.Lock()
	s.entryStreamChans[ech] = struct{}{}
	s.mu.Unlock()

	clientLen := len(s.entryStreamChans)

	logrus.WithFields(logrus.Fields{
		"from":    r.Name,
		"clients": clientLen,
	}).Infof("[ACTION] [>] Connect > %s", r.Name)

	/*entries, err := database.GetAllServerEntry()
	if err == nil {
		for _, e := range entries {
			es.Send(&pb.EntryStreamResponse{Type: pb.StreamType_SYNC, Entry: s.ServerEntry_DBtoPB(e)})
		}
	}*/

	defer func() {
		s.mu.Lock()
		delete(s.entryStreamChans, ech)
		s.mu.Unlock()
		close(ech)
		logrus.WithFields(logrus.Fields{
			"from":    r.Name,
			"clients": clientLen,
		}).Infof("[ACTION] [x] CLOSED > %s", r.Name)
	}()

	for e := range ech {
		if e.Target != "" && !strings.HasPrefix(e.Target, r.Name) {
			continue
		}

		if e.Type == pb.StreamType_QUIT {
			return nil
		}

		logrus.WithFields(logrus.Fields{
			"type": e.Type,
			"from": r.Name,
		}).Debugf("[ACTION] [<->] %s", e.Type)

		err := es.Send(&e)
		if err != nil {
			logrus.WithError(err).Errorf("[gRPC]")
			return err
		}
	}
	return nil
}
