package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/ssubedir/goriffin/internal/handlers"
	proto "github.com/ssubedir/goriffin/protos/service/protos"
	"google.golang.org/grpc"

	"github.com/gorilla/mux"
)

func main() {

	logg := logrus.New()
	logg.Out = os.Stdout
	logg.SetLevel(logrus.InfoLevel)

	l := log.New(os.Stdout, "goriffin-api", log.LstdFlags)

	err := godotenv.Load()
	if err != nil {
		logg.WithFields(logrus.Fields{"error": err}).Fatal("Error loading .env file")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	apicPort := os.Getenv("API_PORT")

	conn, err := grpc.Dial("localhost:"+grpcPort, grpc.WithInsecure()) // do not use WithInsecure in production
	if err != nil {
		panic(err)
	}

	logg.Info("Listening for Goriffin background grpc service on port: ", grpcPort)

	defer conn.Close()

	// create client
	sc := proto.NewServiceClient(conn)

	sh := handlers.NewServices(sc)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	// get
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", sh.Index)
	getR.HandleFunc("/services", sh.GetServices)
	getR.HandleFunc("/service", sh.GetService)
	getR.HandleFunc("/status", sh.GetStatus)

	// post
	getP := sm.Methods(http.MethodPost).Subrouter()
	getP.HandleFunc("/service", sh.AddService)

	// delete
	getD := sm.Methods(http.MethodDelete).Subrouter()
	getD.HandleFunc("/service", sh.RemoveService)

	// create a new server
	s := http.Server{
		Addr:         "localhost:" + apicPort, // configure the bind address
		Handler:      sm,                      // set the default handler
		ErrorLog:     l,                       // set the logger for the server
		ReadTimeout:  5 * time.Second,         // max time to read request from the client
		WriteTimeout: 10 * time.Second,        // max time to write response to the client
		IdleTimeout:  120 * time.Second,       // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logg.Info("Starting Api service on port ", apicPort)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	logg.Info("Got signal:", sig, " Exiting service")

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
