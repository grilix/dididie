package dieded

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.Methods("POST").Path("/profiles").Handler(httptransport.NewServer(
		e.CreateProfileEndpoint,
		decodeCreateProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/profiles/query").Handler(httptransport.NewServer(
		e.QueryProfileEndpoint,
		decodeQueryProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/profiles/{id}/die").Handler(httptransport.NewServer(
		e.DieProfileEndpoint,
		decodeDieProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/profiles/{id}").Handler(httptransport.NewServer(
		e.GetProfileEndpoint,
		decodeGetProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/profiles/{id}").Handler(httptransport.NewServer(
		e.DeleteProfileEndpoint,
		decodeDeleteProfileRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCreateProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var form ProfileForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return nil, err
	}
	return createProfileRequest{
		ProfileForm: form,
	}, nil
}

func decodeGetProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrBadRouting
	}
	return getProfileRequest{ID: idInt}, nil
}

func decodeQueryProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var query Query
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return nil, err
	}
	return queryProfileRequest{
		Query: query,
	}, nil
}

func decodeDieProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrBadRouting
	}

	return dieProfileRequest{
		ID: idInt,
	}, nil
}

func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrBadRouting
	}
	return deleteProfileRequest{ID: idInt}, nil
}

func encodeGetProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getProfileRequest)
	req.URL.Path = fmt.Sprintf("/profiles/%d", r.ID)
	return encodeRequest(ctx, req, request)
}

func encodeQueryProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/profiles/query"
	return encodeRequest(ctx, req, request)
}

func encodeDieProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(dieProfileRequest)
	req.URL.Path = fmt.Sprintf("/profiles/%d/die", r.ID)
	return encodeRequest(ctx, req, request)
}

func encodeDeleteProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(deleteProfileRequest)
	req.URL.Path = fmt.Sprintf("/profiles/%d", r.ID)
	return encodeRequest(ctx, req, request)
}

func decodeGetProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeQueryProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response queryProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDieProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response queryProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	// case ErrAlreadyExists, ErrInconsistentIDs:
	// 	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
