package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/ssubedir/goriffin/internal/data"
	"github.com/ssubedir/goriffin/internal/heartbeat"
	"github.com/ssubedir/goriffin/internal/queue"
	"github.com/ssubedir/goriffin/protos/server"
	proto "github.com/ssubedir/goriffin/protos/service/protos"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logg := logrus.New()
	logg.Out = os.Stdout
	logg.SetLevel(logrus.InfoLevel)

	err := godotenv.Load()
	if err != nil {
		logg.WithFields(logrus.Fields{"error": err}).Info("Error loading .env file")
	}

	logg.Info("Starting Goriffin background worker service")

	workers, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil {
		logg.Panic("Error parsing WORKER_COUNT from .env")
	}

	// grpc
	gs := grpc.NewServer()
	c := server.NewService()
	proto.RegisterServiceServer(gs, c)
	reflection.Register(gs)

	grpcPort := os.Getenv("GRPC_PORT")

	go func() {
		// create a TCP socket for inbound server connections
		l, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			os.Exit(1)
		}

		logg.Info("Starting Goriffin-grpc service on port ", grpcPort)

		// listen for requests
		gs.Serve(l)
	}()

	q := queue.NewQueue(workers)
	q.Start()
	defer q.Stop()

	// DB
	Connections := data.DBConnection()
	mongo := Connections["mongodb"].(data.Mongo)
	mongoConn := mongo.Connect(logg)
	defer mongoConn.Disconnect()

	// Mongo Setup
	goriffin := mongoConn.MONGO.Database("goriffin")
	serviceCollection := goriffin.Collection("service")

	// Fetch collections pass it to queue
	var sHTTP []heartbeat.HTTPService

	sigs := make(chan os.Signal, 1)
	sigdone := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Signals
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		sigdone <- true
	}()

	go func() {
		t := time.NewTicker(time.Second * time.Duration(10))
		defer t.Stop()
		// main loop
		for {
			go func() {

				http, errHTTP := serviceCollection.Find(context.TODO(), bson.M{"stype": "http"})

				if errHTTP != nil {
					q.Stop()
					logg.Panic("Error fetching services")
				}
				if errHTTP = http.All(context.TODO(), &sHTTP); errHTTP != nil {
					q.Stop()
					logg.Panic("Error seralizing http services")
				}

				for _, serv := range sHTTP {

					if heartbeat.FreqCheck(serv.Lastupdate, serv.Frequency) { // check service frequency
						// run task
						go q.Submit(&heartbeat.HTTPHeartbeat{
							Service: serv,
							Logger:  logg,
							Mongo:   mongo.Connect(logg),
						})
					}
				}

			}()
			<-t.C
		}
	}()

	<-sigdone
	logg.Info("Signal received, gracefully stoping Goriffin background Service")
	gs.GracefulStop()
}
