package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next ClarifaiApiService) ClarifaiApiService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	ClarifaiApiService
}

func (mw logmw) PostImage(req PostImageRequest) (resp PostImageResponse, err error) {
	// TODO(madadam): Send to datadog events from here,
	// or write to a file and have something ship logs from the file?
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "PostImage",
			"uri", req.Uri,
			"objectid", resp.ObjectId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	resp, err = mw.ClarifaiApiService.PostImage(req)
	return
}
