package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next ClarifaiAPIService) ClarifaiAPIService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	ClarifaiAPIService
}

func (mw logmw) PostImage(req PostImageRequest) (resp PostImageResponse, err error) {
	// TODO(madadam): Send to datadog events from here,
	// or write to a file and have something ship logs from the file?
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "PostImage",
			"uri", req.URI,
			"objectid", resp.ObjectID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	resp, err = mw.ClarifaiAPIService.PostImage(req)
	return
}

func (mw logmw) GetModels(req GetModelsRequest) (resp GetModelsResponse, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetModels",
			"models", resp.Models,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	resp, err = mw.ClarifaiAPIService.GetModels(req)
	return
}
