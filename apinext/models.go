package main

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ModelInfo struct {
	Name         string   `json:"name"`
	SupportedOps []string `json:"supported_ops"`
}

type GetModelsRequest struct{} // Needed if there's no playload?

type GetModelsResponse struct {
	Models []ModelInfo `json:"models"`
	Err    string      `json:"err,omitempty"`
}

// Aaagh, lack of generics hurts.
func removeDuplicates(s []string) []string {
	var result []string
	seen := map[string]struct{}{}
	for _, v := range s {
		if _, found := seen[v]; !found {
			result = append(result, v)
			seen[v] = struct{}{}
		}
	}
	return result
}

var getModelzResponse = func() (*goquery.Document, error) {
	// Note how we use a function object for dependency injection, see models_test.go.
	// TODO(madadam): base_url as flag/param.
	return goquery.NewDocument("https://api.clarifai.com/v1/modelz")
}

func getModelsFromModelz() (GetModelsResponse, error) {
	var response = GetModelsResponse{
		Models: []ModelInfo{},
		Err:    "",
	}
	doc, err := getModelzResponse()
	if err != nil {
		response.Err = "Error getting model info"
		return response, err
	}
	jsonish := doc.Find("pre").First().Text()
	// It's actually a printed python dict with single quotes, so it's not valid json.  Fix:
	jsonish = strings.Replace(jsonish, "'", "\"", -1)
	var modelzInfo map[string]interface{}
	err = json.Unmarshal([]byte(jsonish), &modelzInfo)
	if err != nil {
		response.Err = "Error getting model info"
		return response, err
	}

	var modelmap = make(map[string][]string)
	for k := range modelzInfo {
		parts := strings.Split(k, ":")
		if _, found := modelmap[parts[0]]; !found {
			modelmap[parts[0]] = make([]string, 0)
		}
		modelmap[parts[0]] = append(modelmap[parts[0]], parts[1])
	}

	models := make([]ModelInfo, 0)
	for name, ops := range modelmap {
		ops = removeDuplicates(ops)
		models = append(models, ModelInfo{name, ops})
	}

	response.Models = models
	return response, nil
}
