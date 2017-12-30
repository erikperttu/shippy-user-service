package main

import (
	"errors"
	"log"

	"context"
	"encoding/json"
	pb "github.com/erikperttu/shippy-user-service/proto/user"
	"github.com/micro/go-micro/broker"
	"golang.org/x/crypto/bcrypt"
)

const topic = "user.created"

type service struct {
	repo         Repository
	tokenService Authable
	PubSub       broker.Broker
}

// From users.pb.go
// Get(context.Context, *User, *Response) error
func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

// From users.pb.go
// GetAll(context.Context, *Request, *Response) error
func (srv *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

// From users.pb.go
// Auth(context.Context, *User, *Token) error
func (srv *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Attemtping to login in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	log.Println(user)

	// Compare our given password against the hashed password
	// stored in the database

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

// From users.pb.go
// Create(context.Context, *User, *Response) error
func (srv *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	// Generate the hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPass)
	//if err := srv.repo.Create(req); err != nil {
	//	return err
	//}
	if err := srv.publishEvent(req); err != nil {
		return err
	}
	res.User = req
	return nil
}

func (srv *service) publishEvent(user *pb.User) error {
	// Marshal to JSON
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Create broker message
	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}

	// Publish the message to the broker
	if err := srv.PubSub.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}
	return nil
}

// From users.pb.go
// ValidateToken(context.Context, *Token, *Token) error
func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	// Decode token
	claims, err := srv.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}
	res.Valid = true
	return nil
}
