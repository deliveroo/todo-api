package selftest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/deliveroo/assert-go"
	"github.com/jackc/fake"
	"github.com/oliveagle/jsonpath"
	"github.com/pkg/errors"
)

// m is shorthand for an untyped map (i.e. a json object).
type m map[string]interface{}

type API struct {
	Username string
	Password string
	Token    string
}

func (a *API) Delete(t *testing.T, path string, body interface{}) *TestResponse {
	t.Helper()
	return a.do(t, http.MethodDelete, path, body)
}

func (a *API) Post(t *testing.T, path string, body interface{}) *TestResponse {
	t.Helper()
	return a.do(t, http.MethodPost, path, body)
}

func (a *API) Put(t *testing.T, path string, body interface{}) *TestResponse {
	t.Helper()
	return a.do(t, http.MethodPut, path, body)
}

func (a *API) Get(t *testing.T, path string) *TestResponse {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url+path, nil)
	assert.Must(t, err)
	if a.Token != "" {
		req.Header.Set("x-todo-token", a.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	assert.Must(t, err)
	return &TestResponse{resp: resp}
}

func (a *API) do(t *testing.T, httpMethod string, path string, body interface{}) *TestResponse {
	var reader io.Reader
	if body != nil {
		json, err := json.Marshal(body)
		assert.Must(t, err)
		reader = bytes.NewBuffer(json)
	}
	req, err := http.NewRequest(httpMethod, url+path, reader)
	assert.Must(t, err)
	if a.Token != "" {
		req.Header.Set("x-todo-token", a.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	assert.Must(t, err)
	return &TestResponse{resp: resp}
}

func fakePassword() string {
	return fake.Password(8, 32, true, true, true)
}

func withAccount(t *testing.T, fn func(*API)) {
	t.Helper()
	api := &API{
		Username: "user-" + time.Now().Format(time.RFC3339Nano),
		Password: fakePassword(),
	}
	account := m{
		"username": api.Username,
		"password": api.Password,
	}
	resp := api.Post(t, "/account", account)
	resp.AssertStatusCode(t, 200)
	resp = api.Post(t, "/account/login", account)
	resp.AssertStatusCode(t, 200)
	api.Token = resp.JSONPathString(t, "token")
	fn(api)
}

type TestResponse struct {
	resp *http.Response
	body []byte
	json interface{}
}

func (r *TestResponse) AssertStatusCode(t *testing.T, want int) {
	t.Helper()
	if r.StatusCode() != want {
		t.Fatalf("unexpected status code, got %d, want %d:\n%s´", r.StatusCode(), want, r.JSONBody(t))
	}
}

func (r *TestResponse) BindBody(t *testing.T, val interface{}) {
	t.Helper()
	r.EnsureReadAll(t)
	if err := json.Unmarshal(r.body, val); err != nil {
		t.Fatalf("%s:\n%s", err, r.body)
	}
}

func (r *TestResponse) EnsureReadAll(t *testing.T) {
	t.Helper()
	if len(r.body) > 0 {
		return
	}
	defer r.resp.Body.Close()
	b, err := ioutil.ReadAll(r.resp.Body)
	assert.Must(t, err)
	r.body = b
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(errors.Wrap(err, "EnsureReadAll Unmarshal"))
	}
	r.json = data
}

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (r *TestResponse) error(t *testing.T) apiError {
	t.Helper()
	r.EnsureReadAll(t)
	var e apiError
	assert.Must(t, json.Unmarshal(r.body, &e))
	return apiError{Error: e.Error}
}

// ErrorCode returns the value under the key 'code' in an Error object. Note
// this is NOT the response status code.
func (r *TestResponse) ErrorCode(t *testing.T) string {
	return r.error(t).Error.Code
}

func (r *TestResponse) ErrorMessage(t *testing.T) string {
	return r.error(t).Error.Message
}

func (r *TestResponse) JSONPathString(t *testing.T, path string) string {
	t.Helper()
	r.EnsureReadAll(t)
	if !strings.HasPrefix(path, "$.") {
		path = "$." + path
	}
	res, err := jsonpath.JsonPathLookup(r.json, path)
	if err != nil {
		t.Fatalf("%s:\n%s´", err, r.JSONBody(t))
	}
	s, ok := res.(string)
	if !ok {
		t.Fatalf("JSON path value %q is not a string", path)
	}
	return s
}

func (r *TestResponse) JSONPath(t *testing.T, path string) interface{} {
	t.Helper()
	r.EnsureReadAll(t)
	if !strings.HasPrefix(path, "$.") {
		path = "$." + path
	}
	res, err := jsonpath.JsonPathLookup(r.json, path)
	if err != nil {
		t.Fatalf("%s:\n%s´", err, r.JSONBody(t))
	}
	return res
}

func (r *TestResponse) JSONPathEqual(t *testing.T, path string, want interface{}) {
	t.Helper()
	r.EnsureReadAll(t)
	if !strings.HasPrefix(path, "$.") {
		path = "$." + path
	}
	res, err := jsonpath.JsonPathLookup(r.json, path)
	if err != nil {
		t.Fatalf("%s:\n%s´", err, r.JSONBody(t))
	}
	assert.Equal(t, res, want)
}

func (r *TestResponse) JSONBody(t *testing.T) string {
	t.Helper()
	r.EnsureReadAll(t)
	return string(r.body)
}

func (r *TestResponse) StatusCode() int {
	return r.resp.StatusCode
}
