package main

import (
	"strings"
	"testing"

	pb "github.com/erikperttu/shippy-user-service/proto/auth"
)

var (
	user = &pb.User{
		Id:    "abc123",
		Email: "test@example.com",
	}
)

// MockRepo empty struct
type MockRepo struct{}

// GetAll get all mocked users
func (repo *MockRepo) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	return users, nil
}

// Get get mocked user by id
func (repo *MockRepo) Get(id string) (*pb.User, error) {
	var user *pb.User
	return user, nil
}

// Create
func (repo *MockRepo) Create(user *pb.User) error {
	return nil
}

// Get mocked user by email
func (repo *MockRepo) GetByEmail(email string) (*pb.User, error) {
	var user *pb.User
	return user, nil
}

// Authable
func newInstance() Authable {
	repo := &MockRepo{}
	return &TokenService{repo}
}

func TestCanCreateToken(t *testing.T) {
	srv := newInstance()
	token, err := srv.Encode(user)
	if err != nil {
		t.Fail()
	}

	if token == "" {
		t.Fail()
	}

	if len(strings.Split(token, ".")) != 3 {
		t.Fail()
	}
}

func TestCanDecodeToken(t *testing.T) {
	srv := newInstance()
	token, err := srv.Encode(user)
	if err != nil {
		t.Fail()
	}
	claims, err := srv.Decode(token)
	if claims.User == nil {
		t.Fail()
	}
	if claims.User.Email != "test@example.com" {
		t.Fail()
	}
}

func TestThrowsErrorIfTokenInvalid(t *testing.T) {
	srv := newInstance()
	_, err := srv.Decode("fail.fail.fail")
	if err == nil {
		t.Fail()
	}
}
