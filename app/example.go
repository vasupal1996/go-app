//go:generate $GOPATH/bin/mockgen -destination=../mock/mock_example.go -package=mock go-app/app Example

package app

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Example defines methods of example service to be implemented
type Example interface {
	SayHello(bool) (string, error)
	SaveHello(*SaveHelloOpts) (*SaveHelloResp, error)
}

// ExampleOpts contains arguments to be accepted for new instance of example service
type ExampleOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// ExampleImpl implements example service
type ExampleImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitExample returns initializes example service
func InitExample(opts *ExampleOpts) Example {
	e := &ExampleImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return e
}

// SayHello returns hello string and error object (if any error)
func (e *ExampleImpl) SayHello(wantErr bool) (string, error) {
	if wantErr {
		return "", errors.New("hello error")
	}
	return "Hello", nil
}

// SaveHelloOpts stores which hello message to save
type SaveHelloOpts struct {
	Name string `json:"name"`
}

// SaveHelloResp returns resp after saving hello in db
type SaveHelloResp struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

// SaveHello save a hello message in the database
func (e *ExampleImpl) SaveHello(opts *SaveHelloOpts) (*SaveHelloResp, error) {
	res, err := e.DB.Collection("hello").InsertOne(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	out := SaveHelloResp{
		Name: opts.Name,
		ID:   res.InsertedID.(primitive.ObjectID),
	}
	return &out, nil
}
