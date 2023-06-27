package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExpectedRequest struct {
	Url    string      `json:"url"`
	Expect interface{} `json:"expect"`
}

func GetTestHttpServer(responses map[string][]byte, expects ...ExpectedRequest) (server *httptest.Server, after func()) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqUrl := req.URL.String()
		if data, ok := responses[reqUrl]; ok {
			//check requests
			for _, expectation := range expects {
				if expectation.Url == reqUrl {
					expectData, _ := json.Marshal(expectation.Expect)
					actualData, _ := io.ReadAll(req.Body)
					if string(expectData) != string(actualData) {
						res.WriteHeader(500)
						msg := fmt.Sprintf("expected actual request data:\n%s \nto equal expected data: %s", string(actualData), string(expectData))
						d, _ := json.Marshal(map[string]string{
							"errorMessage": msg,
						})
						res.Write(d)
						return
					}
				}
			}

			res.Header().Add("Content-Type", "application/json")
			res.WriteHeader(200)
			res.Write(data)
			return
		}

		res.WriteHeader(500)
		msg := fmt.Sprintf("unable to find mock request mapping: %s", reqUrl)
		fmt.Println(msg)
		d, _ := json.Marshal(map[string]string{
			"errorMessage": msg,
		})
		res.Write(d)
	}))

	return testServer, func() {
		testServer.Close()
	}
}

type Request struct {
	Method       string
	Status       int
	Url          string
	Result       []byte
	ExpectInput  interface{}
	ValidateCall func(t *testing.T, req *http.Request)
}

func GetTestHttpServerWithRequests(t *testing.T, reqs []Request) (server *httptest.Server) {
	mReqs := make(map[string]Request)
	for _, req := range reqs {
		mReqs[fmt.Sprintf("%s__%s", req.Method, req.Url)] = req
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqUrl := req.URL.String()
		key := fmt.Sprintf("%s__%s", req.Method, reqUrl)
		r, ok := mReqs[key]
		if !ok {
			res.WriteHeader(500)
			msg := fmt.Sprintf("unable to find mock request mapping: %s", key)
			fmt.Println(msg)
			d, _ := json.Marshal(map[string]string{
				"errorMessage": msg,
			})
			res.Write(d)
			return
		}

		if r.ValidateCall != nil {
			r.ValidateCall(t, req)
		}

		//check requests
		if r.ExpectInput != nil {
			if expected, ok := r.ExpectInput.([]byte); ok {
				var isValid bool
				actualData, _ := io.ReadAll(req.Body)
				defer req.Body.Close()

				switch string(expected[0]) {
				case "{", "[":
					// Assert JSON
					isValid = assert.JSONEq(t, string(expected), string(actualData))
				default:
					isValid = assert.Equal(t, string(expected), string(actualData))
				}

				if !isValid {
					res.WriteHeader(500)
					msg := fmt.Sprintf("expected actual request data:\n%s \nto equal expected data: %s", string(actualData), string(expected))
					d, _ := json.Marshal(map[string]string{
						"errorMessage": msg,
					})
					res.Write(d)
					return
				}
			} else {
				// is not byte array, this is error
				res.WriteHeader(500)
				d, _ := json.Marshal(map[string]string{
					"errorMessage": "expected input (ExpectInput) must be byte array containing JSON",
				})
				res.Write(d)
				return
			}
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(r.Status)
		res.Write(r.Result)
	}))

	t.Cleanup(func() {
		testServer.Close()
	})

	return testServer
}
