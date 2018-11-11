package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

//dec-endpoint-env
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		nanosecond("upper endpoint")
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}
//dec-endpoint-env
func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}
//dec
func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	nanosecond("upper dec")
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
//dec
func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
//env
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	nanosecond("upper - count enc")
	return json.NewEncoder(w).Encode(response)
}
//dec-request
type uppercaseRequest struct {
	S string `json:"s"`
}
//env-response
type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}
//dec-request
type countRequest struct {
	S string `json:"s"`
}
//env-response
type countResponse struct {
	V int `json:"v"`
}
