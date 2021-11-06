package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var httpClient = &http.Client{}

type BasicAuth struct {
	Username string
	Password string
}

func (a BasicAuth) IsZero() bool {
	return a.Username == "" || a.Password == ""
}

type Fetch struct {
	method    string
	url       string
	query     url.Values
	body      io.Reader
	header    http.Header
	basicAuth BasicAuth
	Errors    []error
}

func New() *Fetch {
	return &Fetch{
		method: "GET",
		url:    "",
		query:  url.Values{},
		body:   nil,
		header: http.Header{},
	}
}

func (f *Fetch) Get(url string) *Fetch {
	f.method = "GET"
	f.url = url

	return f
}

func (f *Fetch) Post(url string) *Fetch {
	f.method = "POST"
	f.url = url

	return f
}

func (f *Fetch) Put(url string) *Fetch {
	f.method = "PUT"
	f.url = url

	return f
}

func (f *Fetch) Patch(url string) *Fetch {
	f.method = "PATCH"
	f.url = url

	return f
}

func (f *Fetch) SetQuery(key, value string) *Fetch {
	f.query.Set(key, value)

	return f
}

func (f *Fetch) SetQueryN(kv map[string]string) *Fetch {
	for k, v := range kv {
		f.query.Set(k, v)
	}

	return f
}

func (f *Fetch) AddQuery(key, value string) *Fetch {
	f.query.Add(key, value)
	return f
}

// WithQuery set and overrides current query parameters.
func (f *Fetch) WithQuery(q url.Values) *Fetch {
	f.query = q

	return f
}

// WithHeader overrides existing Header
func (f *Fetch) WithHeader(h http.Header) *Fetch {
	f.header = h
	return f
}

func (f *Fetch) SetHeader(k, v string) *Fetch {
	f.header.Set(k, v)

	return f
}

func (f *Fetch) SetHeaderN(h map[string]string) *Fetch {
	for k, v := range h {
		f.header.Set(k, v)
	}

	return f
}

func (f *Fetch) SetBearerAuth(key string) *Fetch {
	f.header.Set("Authorization", "Bearer "+key)

	return f
}

func (f *Fetch) SetBasicAuth(username, password string) *Fetch {
	f.basicAuth = BasicAuth{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	return f
}

func (f *Fetch) AcceptLang(v string) *Fetch {
	f.header.Set("Accept-Language", v)

	return f
}

func (f *Fetch) Send(body io.Reader) *Fetch {
	f.body = body
	return f
}

func (f *Fetch) StreamJSON(body io.Reader) *Fetch {
	f.header.Add("Content-Type", ContentJSON)
	f.body = body

	return f
}

func (f *Fetch) SendJSONBlob(b []byte) *Fetch {
	return f.StreamJSON(bytes.NewReader(b))
}

func (f *Fetch) SendJSON(v interface{}) *Fetch {
	d, err := json.Marshal(v)
	if err != nil {
		f.Errors = append(f.Errors, err)

		return f
	}

	return f.StreamJSON(bytes.NewReader(d))
}

// End perform the actual HTTP call.
// To get the response body for further processing
// you can use ioutil.ReadAll
// or directly forward raw io to client using
// echo's Stream() method.
func (f *Fetch) End() (*http.Response, []error) {
	if f.Errors != nil {
		return nil, f.Errors
	}

	if len(f.query) != 0 {
		f.url = fmt.Sprintf("%s?%s", f.url, f.query.Encode())
	}

	req, err := http.NewRequest(f.method, f.url, f.body)
	if err != nil {
		f.Errors = append(f.Errors, err)
		return nil, f.Errors
	}

	req.Header = f.header

	resp, err := httpClient.Do(req)
	if err != nil {
		f.Errors = append(f.Errors, err)
		return nil, f.Errors
	}

	return resp, nil
}

// EndBlob reads response body and returns as a slice of bytes
func (f *Fetch) EndBlob() (Response, []error) {
	resp, errs := f.End()
	if errs != nil {
		return Response{}, f.Errors
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		f.Errors = append(f.Errors, err)
		return Response{}, f.Errors
	}

	return Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       b,
	}, nil
}

// EndJSON decodes response body to the specified data structure.
// It also returned the original response body as bytes.
func (f *Fetch) EndJSON(v interface{}) (Response, []error) {
	resp, errs := f.EndBlob()
	if errs != nil {
		return Response{}, f.Errors
	}

	err := json.Unmarshal(resp.Body, v)
	if err != nil {
		f.Errors = append(f.Errors, err)
		return Response{}, f.Errors
	}

	return resp, nil
}
