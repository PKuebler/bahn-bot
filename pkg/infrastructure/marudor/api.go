package marudor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
)

// APIClient for api requests
type APIClient struct {
	EndpointURL *url.URL
	UserAgent   string
	httpClient  *http.Client
	log         *logrus.Entry

	HafasService *HafasService
}

// NewAPIClient return a client with all services
func NewAPIClient(endpoint string, httpClient *http.Client, log *logrus.Entry) (*APIClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	endpointURL, err := url.Parse(fmt.Sprintf("%s/", endpoint))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Trace("init client..")

	c := &APIClient{
		httpClient: httpClient,
		log:        log,

		EndpointURL: endpointURL,
		UserAgent:   "pkuebler/marudor-telegram-bot",
	}

	log.Trace("add services..")
	c.HafasService = &HafasService{client: c}

	return c, nil
}

func (c *APIClient) newAPIRequest(method string, path string, query string, body interface{}) (*http.Request, error) {
	return c.newRequest(c.EndpointURL, method, path, query, body)
}

func (c *APIClient) newRequest(baseURL *url.URL, method string, path string, query string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path, RawQuery: query}
	u := baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			c.log.Fatal(err)
			return nil, err
		}
		c.log.Trace(buf)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		c.log.Fatal(err)
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *APIClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		c.log.Trace(err)
	}
	c.log.Trace(req.URL.String())
	c.log.Trace(string(requestDump))

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res, errors.New(http.StatusText(res.StatusCode))
	}

	err = json.NewDecoder(res.Body).Decode(v)
	return res, err
}
