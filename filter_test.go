package cloudflare

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pageOpts = PaginationOptions{
	PerPage: 25,
	Page:    1,
}

func TestFilter(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expected method 'GET', got %s", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, `{
			"result": {
				"id": "b7ff25282d394be7b945e23c7106ce8a",
				"paused": false,
				"description": "Login from office",
				"expression": "ip.src eq 127.0.0.1"
			},
			"success": true,
			"errors": null,
			"messages": null
		}
		`)
	}

	mux.HandleFunc("/zones/d56084adb405e0b7e32c52321bf07be6/filters/b7ff25282d394be7b945e23c7106ce8a", handler)
	want := Filter{
		ID:          "b7ff25282d394be7b945e23c7106ce8a",
		Paused:      false,
		Description: "Login from office",
		Expression:  "ip.src eq 127.0.0.1",
	}

	actual, err := client.Filter("d56084adb405e0b7e32c52321bf07be6", "b7ff25282d394be7b945e23c7106ce8a")

	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

func TestFilters(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Expected method 'GET', got %s", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, `{
			"result": [
					{
						"id": "b7ff25282d394be7b945e23c7106ce8a",
						"paused": false,
						"description": "Login from office",
						"expression": "ip.src eq 93.184.216.0 and (http.request.uri.path ~ \"^.*/wp-login.php$\" or http.request.uri.path ~ \"^.*/xmlrpc.php$\")"
					},
					{
						"id": "c218c536b2bd406f958f278cf0fa8c0f",
						"paused": false,
						"description": "Login",
						"expression": "(http.request.uri.path ~ \"^.*/wp-login.php$\" or http.request.uri.path ~ \"^.*/xmlrpc.php$\")"
					},
					{
						"id": "f2a64520581a4209aab12187a0081364",
						"paused": false,
						"description": "not /api",
						"expression": "not http.request.uri.path matches \"^/api/.*$\""
			}, {
						"id": "14217d7bd5ab435e84b1bd468bf4fb9f",
						"paused": false,
						"description": "/api",
						"expression": "http.request.uri.path matches \"^/api/.*$\""
			}, {
						"id": "60ee852f9cbb4802978d15600c7f3110",
						"paused": false,
						"expression": "ip.src eq 93.184.216.0"
			} ],
				"success": true,
				"errors": null,
				"messages": null,
				"result_info": {
					"page": 1,
					"per_page": 25,
					"count": 5,
					"total_count": 5,
					"total_pages": 1
			} }
		`)
	}

	mux.HandleFunc("/zones/d56084adb405e0b7e32c52321bf07be6/filters", handler)
	want := []Filter{
		Filter{
			ID:          "b7ff25282d394be7b945e23c7106ce8a",
			Paused:      false,
			Description: "Login from office",
			Expression:  "ip.src eq 93.184.216.0 and (http.request.uri.path ~ \"^.*/wp-login.php$\" or http.request.uri.path ~ \"^.*/xmlrpc.php$\")",
		},
		Filter{
			ID:          "c218c536b2bd406f958f278cf0fa8c0f",
			Paused:      false,
			Description: "Login",
			Expression:  "(http.request.uri.path ~ \"^.*/wp-login.php$\" or http.request.uri.path ~ \"^.*/xmlrpc.php$\")",
		},
		Filter{
			ID:          "f2a64520581a4209aab12187a0081364",
			Paused:      false,
			Description: "not /api",
			Expression:  "not http.request.uri.path matches \"^/api/.*$\"",
		},
		Filter{
			ID:          "14217d7bd5ab435e84b1bd468bf4fb9f",
			Paused:      false,
			Description: "/api",
			Expression:  "http.request.uri.path matches \"^/api/.*$\"",
		}, Filter{
			ID:         "60ee852f9cbb4802978d15600c7f3110",
			Paused:     false,
			Expression: "ip.src eq 93.184.216.0",
		},
	}

	actual, err := client.Filters("d56084adb405e0b7e32c52321bf07be6", pageOpts)

	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}

