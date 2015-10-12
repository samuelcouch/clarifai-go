package main

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
)

func TestRootHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	ctx := context.Background()
	var service ClarifaiApiService
	service = clarifaiApiService{}
	router := makeRouter(ctx, service)

	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("Didn't get %v from /, got %v", http.StatusOK, recorder.Code)
	}
}

func TestPostImage(t *testing.T) {
	expectedUri := "http://foo.com/bar"
	data := map[string]interface{}{
		"uri": expectedUri,
	}
	s, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/images", bytes.NewBuffer(s))
	recorder := httptest.NewRecorder()

	ctx := context.Background()
	var service ClarifaiApiService
	service = clarifaiApiService{}
	router := makeRouter(ctx, service)

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Didn't get %v from /images, got %v", http.StatusOK, recorder.Code)
	}

	var responseData map[string]interface{}
	responseBytes, _ := ioutil.ReadAll(recorder.Body)
	_ = json.Unmarshal([]byte(responseBytes), &responseData)
	// FIXME validation helpers, check key exists.
	if responseData["uri"] != expectedUri {
		t.Errorf("Didn't get expected response['uri'] %v != %v", responseData["uri"], expectedUri)
	}
}
