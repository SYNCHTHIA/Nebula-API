package server

import (
	"golang.org/x/net/context"

	"github.com/sirupsen/logrus"
	pb "gitlab.com/Startail/Nebula-API/nebulapb"
)

func (s *grpcServer) QuitStream(ctx context.Context, e *pb.QuitStreamRequest) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for c := range s.entryStreamChans {
		c <- pb.EntryStreamResponse{Type: pb.StreamType_QUIT}
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
		if e.Type == pb.StreamType_QUIT {
			return nil
		}

		logrus.WithFields(logrus.Fields{
			"type":   e.Type,
			"from":   r.Name,
			"target": e.Target,
		}).Infof("[ACTION] [<->] %s > %s", e.Type, e.Target)

		err := es.Send(&e)
		if err != nil {
			logrus.WithError(err).Errorf("[gRPC]")
			return err
		}
	}
	return nil
}
