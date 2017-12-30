package main

import (
	"fmt"
	"log"

	pb "github.com/erikperttu/shippy-user-service/proto/user"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/registry/mdns"
)

func main() {
	// Creates a database connection and handles the closing before exit
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	// Create a new micro service
	srv := micro.NewService(
		// Must match the protobuf definition
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Init will parse the cmd line flags
	srv.Init()

	// Get instance of the broker using our defaults
	pubSub := srv.Server().Options().Broker

	// Register
	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, pubSub})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
