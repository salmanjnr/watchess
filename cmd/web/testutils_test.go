package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"watchess.org/watchess/pkg/models/mock"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache("./../../ui/html")
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		tournaments:   &mock.TournamentModel{},
		templateCache: templateCache,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
