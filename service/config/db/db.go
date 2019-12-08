package db

import (
	"log"

	"github.com/go-bongo/bongo"
)

var client *bongo.Connection

// Connection ...
var Connection *bongo.Collection

func init() {
	var err error
	clientConfig := &bongo.Config{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "usersTest",
	}
	client, err = bongo.Connect(clientConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connection made successfully")

	Connection = client.Collection("users")
}
