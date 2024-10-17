package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name       string
	url        string
	method     string
	params     []postData
	statusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"sa", "/make-reservation", "GET", []postData{}, http.StatusOK},

	{"sa-post", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"mr-post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Gaurav"},
		{key: "last_name", value: "Kumar"},
		{key: "email", value: "kXKQ8@example.com"},
		{key: "phone", value: "1234567890"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.statusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.statusCode, resp.StatusCode)
			}
		} else {

			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.statusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.statusCode, resp.StatusCode)
			}
		}
	}

}
