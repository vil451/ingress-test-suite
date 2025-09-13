package client_http

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"ingress-test-suite/internal/consts"
	ds "ingress-test-suite/internal/datastruct"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type HTTPTester struct {
	logger Logger
}

func NewHttpTester(logger Logger) *HTTPTester {
	return &HTTPTester{
		logger: logger,
	}
}

func (t *HTTPTester) Test(entry *ds.IngressTestEntry) *ds.TestResult {
	result := &ds.TestResult{
		Host: entry.Host,
		Path: entry.Path,
	}

	url := fmt.Sprintf("http://%s:%d%s", entry.Host, entry.ExtPort, entry.Path)
	t.logger.Printf(consts.MessageRequestURL, url)

	resp, err := http.Get(url)

	if err != nil {
		result.Success = false
		result.Error = errors.Wrap(err, consts.ErrHttpRequestFailed.Error())
		return result
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			t.logger.Printf(errors.Wrap(cerr, consts.ErrFailedCloseResponseBody.Error()).Error())
		}
	}()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode == entry.ExpectedStatus
	return result
}
