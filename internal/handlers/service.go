package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	proto "github.com/ssubedir/goriffin/protos/service/protos"

	"github.com/sirupsen/logrus"
	"github.com/ssubedir/goriffin/internal/data"
	"github.com/ssubedir/goriffin/internal/heartbeat"
	"go.mongodb.org/mongo-driver/bson"
)

type Services struct {
	logger *logrus.Logger
	sc     proto.ServiceClient
}

func NewServices(c proto.ServiceClient) *Services {

	logg := logrus.New()
	logg.Out = os.Stdout
	logg.SetLevel(logrus.InfoLevel)

	return &Services{
		logger: logg,
		sc:     c,
	}

}

func (s *Services) GetServices(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	Connections := data.DBConnection()
	mongo := Connections["mongodb"].(data.Mongo)
	mongoConn := mongo.Connect(s.logger)
	defer mongoConn.Disconnect()

	// Mongo Setup
	goriffin := mongoConn.MONGO.Database("goriffin")
	serviceCollection := goriffin.Collection("service")
	var SERV []heartbeat.HTTPService

	serv, err := serviceCollection.Find(context.TODO(), bson.M{"stype": "http"})

	if err != nil {
		s.logger.Warn("Error fetching http services")
	}
	if err = serv.All(context.TODO(), &SERV); err != nil {
		s.logger.Warn("Error seralizing http services")
	}

	ToJSON(SERV, w)
}

func (s *Services) AddService(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	Connections := data.DBConnection()
	mongo := Connections["mongodb"].(data.Mongo)
	mongoConn := mongo.Connect(s.logger)
	defer mongoConn.Disconnect()

	// Mongo Setup
	goriffin := mongoConn.MONGO.Database("goriffin")
	serviceCollection := goriffin.Collection("service")

	s.logger.Info("Adding New Service")
	var service HTTP
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Enter json data")
	}

	json.Unmarshal(reqBody, &service)

	serviceCollection.InsertOne(context.TODO(), bson.M{"host": service.Host, "name": service.Name, "stype": service.Stype, "accepted_status": service.AcceptedStatus, "frequency": service.Frequency, "request_method": service.RequestMethod, "request_payload": service.RequestPayload, "request_headers": service.RequestHeaders})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)

}

func (s *Services) RemoveService(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	Connections := data.DBConnection()
	mongo := Connections["mongodb"].(data.Mongo)
	mongoConn := mongo.Connect(s.logger)
	defer mongoConn.Disconnect()

	// Mongo Setup
	goriffin := mongoConn.MONGO.Database("goriffin")
	serviceCollection := goriffin.Collection("service")

	s.logger.Info("Removing Service")
	var service HTTP
	reqBody, err := ioutil.ReadAll(r.Body)
	w.WriteHeader(http.StatusCreated)

	if err != nil {
		ToJSON(&Response{false, time.Now().Local().String()}, w)
	}

	json.Unmarshal(reqBody, &service)

	_, err = serviceCollection.DeleteOne(context.TODO(), bson.M{"host": service.Host, "stype": service.Stype})
	if err != nil {
		ToJSON(&Response{false, time.Now().Local().String()}, w)
	}
	ToJSON(&Response{true, time.Now().Local().String()}, w)
}

func (s *Services) GetService(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	Connections := data.DBConnection()
	mongo := Connections["mongodb"].(data.Mongo)
	mongoConn := mongo.Connect(s.logger)
	defer mongoConn.Disconnect()

	// Mongo Setup
	goriffin := mongoConn.MONGO.Database("goriffin")
	serviceCollection := goriffin.Collection("service")
	var SERV []heartbeat.HTTPService
	var service HTTP
	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &service)
	serv, err := serviceCollection.Find(context.TODO(), bson.M{"host": service.Host, "stype": service.Stype})

	if err != nil {
		s.logger.Warn("Error fetching http services")
	}
	if err = serv.All(context.TODO(), &SERV); err != nil {
		s.logger.Warn("Error seralizing http services")
	}

	ToJSON(SERV, w)

}

func (s *Services) GetStatus(w http.ResponseWriter, r *http.Request) {
	resp, _ := s.sc.HeartBeatStatus(context.Background(), &proto.StatusRequest{})
	w.Header().Add("Content-Type", "application/json")

	if resp != nil {
		ToJSON(resp, w)
	} else {
		ToJSON(&Response{false, time.Now().Local().String()}, w)
	}
}

type Response struct {
	Status bool
	Time   string
}

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

type HTTP struct {
	Host           string       ` json:"host"`
	Name           string       ` json:"name"`
	Stype          string       ` json:"stype"`
	AcceptedStatus map[int]bool ` json:"accepted_status"`
	Frequency      float64      ` json:"frequency"`
	RequestMethod  string       ` json:"request_method"`
	RequestPayload string       ` json:"request_payload"`
	RequestHeaders [][2]string  ` json:"request_headers"`
}

func (s *Services) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type status struct {
		Service string
		Status  bool
	}
	ToJSON(&status{"goriffin-api", true}, w)
}

// {
// 	"host":"yahoo.com",
// 	"name":"yahoo",
// 	"stype":"http",
// 	"accepted_status":{
// 		"200":true,
// 		"300":true,
// 		"400":true
// 		},
// 	"frequency":30,
// 	"request_method":"GET",
// 	"request_payload":"",
// 	"request_headers":[["Host","yahoo.com"]]
// }

// {
// 	"host":"yahoo.com",
// 	"stype":"http"
// }
