package heartbeat

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ssubedir/goriffin/internal/data"
)

type HTTPHeartbeat struct {
	ID      string
	Service HTTPService
	Logger  *logrus.Logger
	Mongo   *data.MongoClient
}

func FreqCheck(t time.Time, freq float64) bool {

	if time.Since(t).Seconds() > freq {
		return true
	} else {
		return false
	}
}
