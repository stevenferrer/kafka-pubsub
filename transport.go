package usermanagementsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(endpoints Endpoints) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	mux := http.NewServeMux()

	mux.Handle("/createuser", httptransport.NewServer(
		endpoints.CreateUserEndpoint,
		DecodeHTTPCreateUserRequest,
		EncodeHTTPCreateUserResponse,
		options...,
	))

	return mux
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}

	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string
}

func EncodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func DecodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func EncodeHTTPCreateUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(createUserResponse)

	r := struct {
		Err string
	}{err2str(resp.Err)}

	err := json.NewEncoder(w).Encode(r)
	return err
}

func DecodeHTTPCreateUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}

	var rp struct {
		Err string
	}

	resp := createUserResponse{
		Err: str2err(rp.Err),
	}

	return resp, nil
}
