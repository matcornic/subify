// Fluent HTTP client for Golang. With timeout, retries and exponential back-off support.
package fluent

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/lafikl/backoff"
)

type Request struct {
	header    map[string]string
	method    string
	json      interface{}
	jsonIsSet bool
	url       string
	retry     int
	timeout   time.Duration
	body      io.Reader
	res       *http.Response
	err       error
	backoff   *backoff.ExponentialBackOff
	req       *http.Request
	proxy     string
}

func (f *Request) newClient() *http.Client {
	return &http.Client{Timeout: f.timeout}
}

func (f *Request) newRequest() (req *http.Request, err error) {
	if f.jsonIsSet {
		body, jsonErr := json.Marshal(f.json)
		if jsonErr != nil {
			return nil, jsonErr
		}
		req, err = http.NewRequest(f.method, f.url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(f.method, f.url, f.body)
	}
	return
}

// Set the request URL
// You probably want to use the methods [Post, Get, Patch, Delete, Put]
func (f *Request) Url(url string) *Request {
	f.url = url
	return f
}

// Set the request Method
// You probably want to use the methods [Post, Get, Patch, Delete, Put]
func (f *Request) Method(method string) *Request {
	f.method = method
	return f
}

// This is a shorthand method that calls f.Method with `POST`
// and calls f.Url with the url you give to her
func (f *Request) Post(url string) *Request {
	f.Url(url).Method("POST")
	return f
}

// Same as f.Post but the method is `PUT`
func (f *Request) Put(url string) *Request {
	f.Url(url).Method("PUT")
	return f
}

// Same as f.Post but the method is `PATCH`
func (f *Request) Patch(url string) *Request {
	f.Url(url).Method("PATCH")
	return f
}

// Same as f.Post but the method is `GET`
func (f *Request) Get(url string) *Request {
	f.Url(url).Method("GET")
	return f
}

// Same as f.Post but the method is `DELETE`
func (f *Request) Delete(url string) *Request {
	f.Url(url).Method("DELETE")
	return f
}

// A handy method for sending json without needing to Marshal it yourself
// This method will override whatever you pass to f.Body
// And it sets the content-type to "application/json"
func (f *Request) Json(j interface{}) *Request {
	f.json = j
	f.jsonIsSet = true
	f.SetHeader("Content-type", "application/json")
	return f
}

// Whatever you pass to it will be passed to http.NewRequest
func (f *Request) Body(b io.Reader) *Request {
	f.body = b
	return f
}

// sets the header entries associated with key to the element value.
//
// It replaces any existing values associated with key.
func (f *Request) SetHeader(key, value string) *Request {
	f.header[key] = value
	return f
}

// Timeout specifies a time limit for requests made by this
// Client. The timeout includes connection time, any
// redirects, and reading the response body. The timer remains
// running after Get, Head, Post, or Do return and will
// interrupt reading of the Response.Body.
//
// A Timeout of zero means no timeout.
func (f *Request) Timeout(t time.Duration) *Request {
	f.timeout = t
	return f
}

// The initial interval for the request backoff operation.
//
// the default is 500 * time.Millisecond
//
// For more information http://godoc.org/github.com/cenkalti/backoff#ExponentialBackOff
func (f *Request) InitialInterval(t time.Duration) *Request {
	f.backoff.InitialInterval = t
	return f
}

// Set the Randomization factor for the backoff.
//
// the default is 0.5
//
// For more information http://godoc.org/github.com/cenkalti/backoff#ExponentialBackOff
func (f *Request) RandomizationFactor(rf float64) *Request {
	f.backoff.RandomizationFactor = rf
	return f
}

// Set the Multiplier for the backoff.
//
// The default is 1.5.
//
// For more information http://godoc.org/github.com/cenkalti/backoff#ExponentialBackOff
func (f *Request) Multiplier(m float64) *Request {
	f.backoff.Multiplier = m
	return f
}

// Set the Max Interval for the backoff.
//
// The default is 60 * time.Second
//
// For more information http://godoc.org/github.com/cenkalti/backoff#ExponentialBackOff
func (f *Request) MaxInterval(mi time.Duration) *Request {
	f.backoff.MaxInterval = mi
	return f
}

// Set the Max Elapsed Time for the backoff.
//
// The default is 15 * time.Minute.
//
// For more information http://godoc.org/github.com/cenkalti/backoff#ExponentialBackOff
func (f *Request) MaxElapsedTime(me time.Duration) *Request {
	f.backoff.MaxElapsedTime = me
	return f
}

// Set how many times to retry if the request
// timedout or the server returned 5xx response.
func (f *Request) Retry(r int) *Request {
	f.retry = r
	return f
}

// Set a HTTP proxy
func (f *Request) Proxy(p string) *Request {
	f.proxy = p
	return f
}

func doReq(f *Request, c *http.Client) error {
	var reqErr error
	f.req, reqErr = f.newRequest()
	if reqErr != nil {
		return reqErr
	}
	for k, v := range f.header {
		f.req.Header.Set(k, v)
	}
	res, err := c.Do(f.req)
	// if there's an error in the request
	// then we just return whatever err we got
	if err != nil {
		f.err = err
		return nil
	}
	if res != nil && res.StatusCode >= 500 && res.StatusCode <= 599 && f.retry > 0 {
		f.retry--
		return errors.New("Server Error")
	}
	if res != nil {
		f.res = res
	}
	return nil
}

func (f *Request) operation(c *http.Client) func() error {
	return func() error {
		return doReq(f, c)
	}
}

func (f *Request) do(c *http.Client) (*http.Response, error) {
	err := doReq(f, c)
	if err != nil {
		op := f.operation(c)
		err = backoff.Retry(op, f.backoff)
		if err != nil {
			return nil, err
		}
	}
	// Check if has operation failed after the retries
	if f.err != nil {
		return nil, f.err
	}
	return f.res, err
}

// It will construct the client and the request, then send it
//
// This function has to be called as the last thing,
// after setting the other properties
func (f *Request) Send() (*http.Response, error) {
	c := *http.DefaultClient
	if f.timeout != 0 {
		nc := f.newClient()
		c = *nc
	}

	if f.proxy != "" {
		proxyUrl, err := url.Parse(f.proxy)

		if err != nil {
			return nil, err
		}

		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	res, err := f.do(&c)
	return res, err
}

// Create a new request
func New() *Request {
	f := &Request{}
	f.header = map[string]string{}
	f.backoff = backoff.NewExponentialBackOff()
	f.err = nil
	return f
}
