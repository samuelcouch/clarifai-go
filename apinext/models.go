package main

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Model struct {
	Name         string   `json:"model"`
	SupportedOps []string `json:"supported_ops"`
}

type GetModelsRequest struct{} // Needed if there's no playload?

type GetModelsResponse struct {
	Models []Model `json:"models"`
	Err    string  `json:"err,omitempty"`
}

// Aaagh, lack of generics hurts.
func removeDuplicates(s []string) []string {
	result := []string{}
	seen := map[string]struct{}{}
	for _, v := range s {
		if _, found := seen[v]; !found {
			result = append(result, v)
			seen[v] = struct{}{}
		}
	}
	return result
}

func getModelsFromModelz() (GetModelsResponse, error) {
	// TODO(madadam): base_url as flag/param.
	doc, err := goquery.NewDocument("https://api.clarifai.com/v1/modelz")
	if err != nil {
		var response = GetModelsResponse{
			Models: []Model{},
			Err:    "Error getting model info",
		}
		return response, err
	}
	jsonish := doc.Find("pre").First().Text()
	// It's actually a printed python dict with single quotes, so it's not valid json.  Fix:
	jsonish = strings.Replace(jsonish, "'", "\"", -1)
	var f interface{}
	err = json.Unmarshal([]byte(jsonish), &f)
	if err != nil {
		var response = GetModelsResponse{
			Models: []Model{},
			Err:    "Error getting model info",
		}
		return response, err
	}
	modelzInfo := f.(map[string]interface{})

	var modelmap = make(map[string][]string)
	for k, _ := range modelzInfo {
		parts := strings.Split(k, ":")
		if _, found := modelmap[parts[0]]; !found {
			modelmap[parts[0]] = make([]string, 0)
		}
		modelmap[parts[0]] = append(modelmap[parts[0]], parts[1])
	}

	models := make([]Model, 0)
	for name, ops := range modelmap {
		ops = removeDuplicates(ops)
		models = append(models, Model{name, ops})
	}

	return GetModelsResponse{
		Models: models,
		Err:    "",
	}, nil
}
