package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDD struct {
	client *mongo.Client
}

type IMongoDB interface {
	Client() *mongo.Client
}

type MongoDBParam struct {
	Host             string
	Database         string
	Port             int
	User             string
	Password         string
	UsedMongoReplica int
}

func NewMongoDB(param MongoDBParam) IMongoDB {
	var mongoURL string

	if param.UsedMongoReplica == 0 {
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin&replicaSet=rs0", param.User, param.Password, param.Host, param.Port, param.Database)
	} else if param.UsedMongoReplica == 2 {
		mongoURL = param.Host
	} else {
		mongoURL = fmt.Sprintf("mongodb+srv://%s:%s@%s", param.User, param.Password, param.Host)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))

	if err != nil {
		panic(err)
	}

	// Auto-reconnect
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				// Periksa koneksi setiap 1 menit
				err := client.Ping(context.Background(), nil)
				if err != nil {
					// Jika koneksi terputus, coba reconnect
					log.Println("Reconnecting to MongoDB...")
					client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
					if err != nil {
						log.Println("Failed to reconnect:", err)
					} else {
						err = client.Ping(context.Background(), nil)
						if err == nil {
							log.Println("Reconnected to MongoDB successfully.")
						} else {
							log.Println("Failed to reconnect:", err)
						}
					}
				}
			}
		}
	}()

	return &MongoDD{
		client: client,
	}
}

func (m *MongoDD) Client() *mongo.Client {
	return m.client
}
