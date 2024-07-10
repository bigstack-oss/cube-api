package http

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	log "go-micro.dev/v5/logger"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Helper struct {
	Client
	Config
}

type Config struct {
	Name   string
	Scheme string
	Host   string

	URL  string
	URLs map[string]string

	Port  int
	Ports map[string]int

	Path  string
	Paths map[string]string

	Method  string
	Methods map[string]string
	Body    string

	CACert string
	BasicAuth
	TLSInsecureSkip bool
	Headers         []Header

	Fetch
}

type Fetch struct {
	Interval time.Duration
	Retry    int
}

type Header struct {
	Key   string
	Value string
}

type BasicAuth struct {
	Username string
	Password string
}

func NewHelper() *Helper {
	h := &Helper{}
	h.SetHttpClient()
	return h
}

func (h *Helper) SetHttpClient() {
	h.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: h.Config.TLSInsecureSkip,
			},
		},
		Timeout: 10 * time.Second,
	}

	if h.Config.Retry == 0 {
		h.Config.Retry = 3
	}
}

func (h *Helper) SetHeadersOnReq(req *http.Request) {
	for _, h := range h.Headers {
		req.Header.Add(h.Key, h.Value)
	}
}

func (h *Helper) Send(req *http.Request) ([]byte, int) {
	trialCount := 0
	statusCode := 0
	respBody := []byte{}

	for {
		if trialCount > h.Config.Retry {
			log.Debugf(
				"reach the http retry limit for URL(%s)",
				req.URL.String(),
			)
			return respBody, statusCode
		}

		resp, err := h.Client.Do(req)
		if err != nil {
			log.Errorf("error details of do http request: %s", err.Error())
			trialCount++
			continue
		}

		if Is5XXCode[resp.StatusCode] {
			log.Debugf("get retryable status code: %d", resp.StatusCode)
			trialCount++
			statusCode = resp.StatusCode
			resp.Body.Close()
			time.Sleep(time.Second * 2)
			continue
		}

		statusCode = resp.StatusCode
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("error details of read resp body: %s", err.Error())
			trialCount++

			continue
		}

		return respBody, statusCode
	}
}

func (h *Helper) GetHeaderVal(req *http.Request, header string) (string, int) {
	statusCode := 0
	trialCount := 0

	for {
		if trialCount > h.Config.Retry {
			log.Errorf(
				"reach the http retry limit for URL(%s)",
				req.URL.String(),
			)
			return "", statusCode
		}

		resp, err := h.Client.Do(req)
		if err != nil {
			log.Errorf("error details of do http request: %s", err.Error())
			trialCount++
			continue
		}

		if Is5XXCode[resp.StatusCode] {
			log.Debugf("get retryable status code: %d", resp.StatusCode)
			trialCount++
			statusCode = resp.StatusCode
			resp.Body.Close()
			time.Sleep(time.Second * 2)
			continue
		}

		statusCode = resp.StatusCode
		headerVal := resp.Header.Get(header)
		resp.Body.Close()
		if err != nil {
			log.Errorf("error details of get header(%s) val: %s", header, err.Error())
			trialCount++
			continue
		}

		return headerVal, statusCode
	}
}
