package dieded

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	CreateProfileEndpoint endpoint.Endpoint
	GetProfileEndpoint    endpoint.Endpoint
	QueryProfileEndpoint  endpoint.Endpoint
	DieProfileEndpoint    endpoint.Endpoint
	DeleteProfileEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateProfileEndpoint: MakeCreateProfileEndpoint(s),
		GetProfileEndpoint:    MakeGetProfileEndpoint(s),
		QueryProfileEndpoint:  MakeQueryProfileEndpoint(s),
		DieProfileEndpoint:    MakeDieProfileEndpoint(s),
		DeleteProfileEndpoint: MakeDeleteProfileEndpoint(s),
	}
}

func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		GetProfileEndpoint:    httptransport.NewClient("GET", tgt, encodeGetProfileRequest, decodeGetProfileResponse, options...).Endpoint(),
		QueryProfileEndpoint:  httptransport.NewClient("POST", tgt, encodeQueryProfileRequest, decodeQueryProfileResponse, options...).Endpoint(),
		DieProfileEndpoint:    httptransport.NewClient("POST", tgt, encodeDieProfileRequest, decodeDieProfileResponse, options...).Endpoint(),
		DeleteProfileEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteProfileRequest, decodeDeleteProfileResponse, options...).Endpoint(),
	}, nil
}

func (e Endpoints) GetProfile(ctx context.Context, id int) (*Profile, error) {
	request := getProfileRequest{ID: id}
	response, err := e.GetProfileEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(getProfileResponse)
	return resp.Profile, resp.Err
}

func (e Endpoints) DeleteProfile(ctx context.Context, id int) error {
	request := deleteProfileRequest{ID: id}
	response, err := e.DeleteProfileEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteProfileResponse)
	return resp.Err
}

func MakeCreateProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createProfileRequest)
		p, e := s.CreateProfile(ctx, req.ProfileForm)
		return createProfileResponse{Profile: p, Err: e}, nil
	}
}

func MakeGetProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProfileRequest)
		p, e := s.GetProfile(ctx, req.ID)
		return getProfileResponse{Profile: p, Err: e}, nil
	}
}

func MakeQueryProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(queryProfileRequest)
		p, e := s.QueryProfile(ctx, req.Query)
		return queryProfileResponse{Profile: p, Err: e}, nil
	}
}

func MakeDieProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dieProfileRequest)
		p, e := s.DieProfile(ctx, req.ID)
		return dieProfileResponse{Profile: p, Err: e}, nil
	}
}

func MakeDeleteProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteProfileRequest)
		e := s.DeleteProfile(ctx, req.ID)
		return deleteProfileResponse{Err: e}, nil
	}
}

type createProfileRequest struct {
	ProfileForm
}

type getProfileRequest struct {
	ID int
}

type createProfileResponse struct {
	Profile *Profile `json:"profile,omitempty"`
	Err     error    `json:"err,omitempty"`
}

type getProfileResponse struct {
	Profile *Profile `json:"profile,omitempty"`
	Err     error    `json:"err,omitempty"`
}

func (r getProfileResponse) error() error { return r.Err }

type queryProfileRequest struct {
	Query
}

type dieProfileRequest struct {
	ID int
}

type queryProfileResponse struct {
	Profile *Profile `json:"profile,omitempty"`
	Err     error    `json:"err,omitempty"`
}

func (r queryProfileResponse) error() error { return r.Err }

type dieProfileResponse struct {
	Profile *Profile `json:"profile,omitempty"`
	Err     error    `json:"err,omitempty"`
}

func (r dieProfileResponse) error() error { return r.Err }

type deleteProfileRequest struct {
	ID int
}

type deleteProfileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteProfileResponse) error() error { return r.Err }
