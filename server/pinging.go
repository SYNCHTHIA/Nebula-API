package server

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Nebula-API/database"
	"gitlab.com/Startail/Nebula-API/util"
)

func (s *grpcServer) pinging() {
	e, err := database.GetAllServerEntry()
	if err != nil {
		logrus.Errorf("[Database] Error %s", err)
	} else {
		for _, v := range e {
			logrus.Debugf("Trying Entry: %s", v.Name)
			go func(data database.ServerData) {
				//s.mu.Lock()
				r, pingErr := util.Ping(data.Address + ":" + fmt.Sprint(data.Port))

				if r != nil && pingErr == nil {
					logrus.Debugf("%s %d: %s", data.Name, data.Port, r)
					data.Status = *r
				} else {
					logrus.Debugf("%s is offline", data.Name)
					data.Status = database.PingResponse{}
				}
				_, updated, pushErr := database.PushServerStatus(data.Name, data.Status)
				//s.mu.Unlock()
				if pushErr != nil {
					return
				}

				if updated > 0 {
					database.PublishServer(s.ServerEntry_DBtoPB(data))
				}
			}(v)
		}
	}
}
