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
	"github.com/go-zoo/bone"
	"github.com/gorilla/mux"
	"github.com/zenazn/goji"
)

// ClarifaiAPIService is the main entry point to the Clarifai API.
type ClarifaiAPIService interface {
	PostImage(PostImageRequest) (PostImageResponse, error)
}

type clarifaiAPIService struct{}

func (clarifaiAPIService) PostImage(request PostImageRequest) (PostImageResponse, error) {
	// FIXME: schema validation. swagger?  or manually?
	if false {
		// FIXME testing errors
		return PostImageResponse{"", "", "", "bad stuff happened"},
			&APIError{500, "boom, error occurred.", "you broke it!"}
	}
	var response = PostImageResponse{
		"Ed1nuqPvcm",
		"2011-08-20T02:06:57.931Z",
		request.URI,
		"",
	}
	return response, nil
}

// ServiceMiddleware is a chainable middleware type.
type ServiceMiddleware func(ClarifaiAPIService) ClarifaiAPIService

// APIError is the package-wide error representation.
type APIError struct {
	HTTPStatus int
	UserMsg    string
	LogMsg     string
}

func (err *APIError) Error() string {
	return fmt.Sprintf("%d %s [%s]", err.HTTPStatus, err.UserMsg, err.LogMsg)
}

// API schema object types.

// PostImageRequest is the structure through which requests are issued.
type PostImageRequest struct {
	URI string `json:"uri"`
}

// PostImageResponse is the response to a PostImageRequest.
type PostImageResponse struct {
	ObjectID  string `json:"objectId"`
	CreatedAt string `json:"createdAt"`
	URI       string `json:"uri"`
	Err       string `json:"err,omitempty"` // errors don't define JSON marshaling
}

// Route defines how to map HTTP endpoints to handlers.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

// Routes is a collection of Route objects.
type Routes []Route

func makeRoutes(ctx context.Context, service ClarifaiAPIService) *Routes {
	postImageHandler := httptransport.NewServer(
		ctx,
		makePostImageEndpoint(service),
		decodePostImageRequest,
		encodeResponse,
	)

	proxy, err := NewProxy("https://api.clarifai.com")
	if err != nil {
		panic("Couldn't create proxy handler.")
	}
	proxyHandler := httptransport.NewServer(
		ctx,
		NewProxyEndpoint(proxy),
		decodeProxyRequest,
		encodeRecordedResponse,
	)

	var routes = Routes{
		Route{
			"Images",
			"POST",
			"/images",
			postImageHandler,
		},
		Route{
			"Proxy",
			"*",
			"/v1/*",
			proxyHandler,
		},
	}
	return &routes
}

func makeGorillaRouter(ctx context.Context, service ClarifaiAPIService) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	routes := makeRoutes(ctx, service)
	for _, route := range *routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	return router
}

func makeGojiRouter(ctx context.Context, service ClarifaiAPIService) http.Handler {
	routes := makeRoutes(ctx, service)
	for _, route := range *routes {
		// Hm... these are unexported, Goji wants to hide them. Need to iterate all types..
		//method := goji.web.httpMethod(route.Method)
		//goji.DefaultMux.handleUntyped(route.Pattern, method, route.Handler)
		switch {
		case route.Method == "DELETE":
			goji.Delete(route.Pattern, route.Handler)
		case route.Method == "GET":
			goji.Get(route.Pattern, route.Handler)
		case route.Method == "POST":
			goji.Post(route.Pattern, route.Handler)
		case route.Method == "PUT":
			goji.Put(route.Pattern, route.Handler)
		case true:
			panic(fmt.Sprintf("error, unknown method: %v", route.Method))
		}
	}

	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello goodbye %q", html.EscapeString(r.URL.Path))
	})

	goji.DefaultMux.Compile()
	return goji.DefaultMux
}

func makeBoneRouter(ctx context.Context, service ClarifaiAPIService) http.Handler {
	mux := bone.New()
	routes := makeRoutes(ctx, service)
	for _, route := range *routes {
		// TODO(madadam): Boo, again, hit an unexported method (register).
		//mux.register(route.Method, route.Pattern, route.Handler)
		switch {
		case route.Method == "DELETE":
			mux.Delete(route.Pattern, route.Handler)
		case route.Method == "GET":
			mux.Get(route.Pattern, route.Handler)
		case route.Method == "POST":
			mux.Post(route.Pattern, route.Handler)
		case route.Method == "PUT":
			mux.Put(route.Pattern, route.Handler)
		case route.Method == "*":
			mux.Handle(route.Pattern, route.Handler)
		case true:
			panic(fmt.Sprintf("error, unknown method: %v", route.Method))
		}
	}

	mux.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	}))

	return mux
}

func main() {
	var (
		listen     = flag.String("listen", ":8080", "HTTP port")
		routerType = flag.String("router", "bone", "Router package name (bone, goji, gorilla)")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	// Redirect stdlib log to gokit's logger.
	// import (stdlog "log")
	//stdlog.SetOutput(log.NewStdlibAdapter(logger))

	ctx := context.Background()
	var service ClarifaiAPIService
	service = clarifaiAPIService{}
	service = loggingMiddleware(logger)(service)

	var router http.Handler
	switch *routerType {
	default:
		panic(fmt.Sprintf("Unknown router type: %v", *routerType))
	case "gorilla":
		router = makeGorillaRouter(ctx, service)
	case "goji":
		router = makeGojiRouter(ctx, service)
	case "bone":
		router = makeBoneRouter(ctx, service)
	}

	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, router))
}

func makePostImageEndpoint(svc ClarifaiAPIService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostImageRequest)
		response, err := svc.PostImage(req)
		if err != nil {
			// FIXME error handling
			return PostImageResponse{"", "", "", err.Error()},
				&APIError{500, "Sorry, an error occurred.", err.Error()}
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
