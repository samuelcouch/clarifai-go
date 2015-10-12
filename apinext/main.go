package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type ClarifaiApiService interface {
	PostImage(PostImageRequest) (PostImageResponse, error)
}

type clarifaiApiService struct{}

func (clarifaiApiService) PostImage(request PostImageRequest) (PostImageResponse, error) {
	// FIXME: schema validation. swagger?  or manually?
	if false {
		// FIXME testing errors
		return PostImageResponse{"", "", "", "bad stuff happened"},
			&ApiError{500, "boom, error occurred.", "you broke it!"}
	}
	var response = PostImageResponse{
		"Ed1nuqPvcm",
		"2011-08-20T02:06:57.931Z",
		request.Uri,
		"",
	}
	return response, nil
}

// Chainable middleware type.
type ServiceMiddleware func(ClarifaiApiService) ClarifaiApiService

type ApiError struct {
	HttpStatus int
	UserMsg    string
	LogMsg     string
}

func (err *ApiError) Error() string {
	return fmt.Sprintf("%d %s [%s]", err.HttpStatus, err.UserMsg, err.LogMsg)
}

// API schema object types.

type PostImageRequest struct {
	Uri string `json:"uri"`
}

type PostImageResponse struct {
	ObjectId  string `json:"objectId"`
	CreatedAt string `json:"createdAt"`
	Uri       string `json:"uri"`
	Err       string `json:"err,omitempty"` // errors don't define JSON marshaling
}

// Routes.

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

func makeRoutes(ctx context.Context, service ClarifaiApiService) *Routes {
	postImageHandler := httptransport.NewServer(
		ctx,
		makePostImageEndpoint(service),
		decodePostImageRequest,
		encodeResponse,
	)

	var routes = Routes{
		Route{
			"Images",
			"POST",
			"/images",
			postImageHandler,
		},
	}
	return &routes
}

func makeRouter(ctx context.Context, service ClarifaiApiService) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	routes := makeRoutes(ctx, service)
	for _, route := range *routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}

	// FIXME
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	return router
}

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP port")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	ctx := context.Background()
	var service ClarifaiApiService
	service = clarifaiApiService{}
	service = loggingMiddleware(logger)(service)

	router := makeRouter(ctx, service)

	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, router))
}

func makePostImageEndpoint(svc ClarifaiApiService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostImageRequest)
		response, err := svc.PostImage(req)
		if err != nil {
			// FIXME error handling
			return PostImageResponse{"", "", "", err.Error()},
				&ApiError{500, "Sorry, an error occurred.", err.Error()}
		}
		return response, err
	}
}

func decodePostImageRequest(r *http.Request) (interface{}, error) {
	var request PostImageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// FIXME: how to get error codes with http status?
	return json.NewEncoder(w).Encode(response)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
