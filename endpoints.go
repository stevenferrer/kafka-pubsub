package usermanagementsvc

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUserEndpoint endpoint.Endpoint
}

type createUserRequest struct {
	Email    string
	Password string
}

type createUserResponse struct {
	Err error
}

func (e Endpoints) CreateUser(ctx context.Context, email, password string) error {
	request := createUserRequest{Email: email, Password: password}
	response, err := e.CreateUserEndpoint(ctx, request)
	if err != nil {
		return err
	}

	r := response.(createUserResponse)
	return r.Err
}

func MakeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		createUserReq := request.(createUserRequest)
		err := s.CreateUser(ctx, createUserReq.Email, createUserReq.Password)
		return createUserResponse{
			Err: err,
		}, nil
	}
}
