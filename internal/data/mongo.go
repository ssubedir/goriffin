package data

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Username string
	Password string
	Cluster  string
}

type MongoClient struct {
	MONGO  *mongo.Client
	Logger *logrus.Logger
}

func (m Mongo) Connect(l *logrus.Logger) *MongoClient {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://"+m.Username+":"+m.Password+"@"+m.Cluster+"/test?retryWrites=true&w=majority",
	))
	if err != nil {
		l.WithFields(logrus.Fields{"error": err}).Fatal("Error connecting to MongoDB")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		l.WithFields(logrus.Fields{"error": err}).Fatal("Error pinging MongoDB")
	} else {
		l.Debug("Connected to MongoDB")
	}

	return &MongoClient{
		MONGO:  client,
		Logger: l,
	}

}

func (m *MongoClient) Disconnect() {
	err := m.MONGO.Disconnect(context.TODO())
	if err != nil {
		m.Logger.WithFields(logrus.Fields{"error": err}).Fatal("Error closing MONGO DB")
	} else {
		m.Logger.Debug("Disconnected from MongoDB")

	}
}
