package heartbeat

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTPService struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Host           string             `bson:"host,omitempty" json:"host"`
	Name           string             `bson:"name,omitempty" json:"name"`
	Stype          string             `bson:"stype,omitempty" json:"stype"`
	AcceptedStatus map[int]bool       `bson:"accepted_status,omitempty" json:"accepted_status"`
	Frequency      float64            `bson:"frequency,omitempty" json:"frequency"`
	RequestMethod  string             `bson:"request_method,omitempty" json:"request_method"`
	RequestPayload string             `bson:"request_payload,omitempty" json:"request_payload"`
	RequestHeaders [][2]string        `bson:"request_headers,omitempty" json:"request_headers"`
	Lastupdate     time.Time          `bson:"last_update,omitempty"`
}

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

// Run - Http
func (h *HTTPHeartbeat) Run() {
	defer h.Mongo.Disconnect()

	// Do task
	h.CheckHeartbeat()

}

// CheckHeartbeat - Http
func (h *HTTPHeartbeat) CheckHeartbeat() {

	// update Db
	client := h.Mongo.MONGO
	goriffin := client.Database(os.Getenv("DB_DBNAME"))
	httpCollection := goriffin.Collection(os.Getenv("HEARTBEAT_HTTP"))
	serviceCollection := goriffin.Collection(os.Getenv("SERVICES_HEARTBEAT"))

	// Setup
	body := strings.NewReader(h.Service.RequestPayload)
	req, err := http.NewRequest(h.Service.RequestMethod, h.Service.Stype+"://"+h.Service.Host, body)
	if err != nil {
		h.Logger.Info(h.Service.Host, " Error creating http request, make sure request data is valid")
		return
	}

	// Add custom headers
	for i := range h.Service.RequestHeaders {
		req.Header.Set(h.Service.RequestHeaders[i][0], h.Service.RequestHeaders[i][1])
	}

	// Send Request
	start := time.Now()
	response, err := netClient.Do(req)
	if err != nil {

		if strings.Contains(err.Error(), "no such host") {
			h.Logger.Info(h.Service.Host, " - No Such host, Check Connection / DNS")
		} else if strings.Contains(err.Error(), "timeout") {
			h.Logger.Info(h.Service.Host, " - Http timed out")
		} else {
			h.Logger.Info(h.Service.Host, " - Error Sending http request")
		}
		serviceCollection.UpdateOne(context.TODO(), bson.M{"_id": h.Service.ID}, bson.D{{"$set", bson.D{{"last_update", time.Now()}, {"status", 0}, {"alive", false}}}})

		return
	}
	stop := time.Since(start).Milliseconds()

	// close body
	defer response.Body.Close()

	h.Logger.Debug(h.Service.Host, " ", response.StatusCode)

	// Acceptable status code
	_, ok := h.Service.AcceptedStatus[response.StatusCode]
	if ok {
		httpCollection.InsertOne(context.TODO(), bson.M{"host": h.Service.Host, "status": response.StatusCode, "rtt": stop, "time": time.Now()})
		serviceCollection.UpdateOne(context.TODO(), bson.M{"_id": h.Service.ID}, bson.D{{"$set", bson.D{{"last_update", time.Now()}, {"status", response.StatusCode}, {"alive", true}}}})
	} else {
		serviceCollection.UpdateOne(context.TODO(), bson.M{"_id": h.Service.ID}, bson.D{{"$set", bson.D{{"last_update", time.Now()}, {"status", response.StatusCode}, {"alive", false}}}})
	}

}
