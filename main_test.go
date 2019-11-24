package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCloudProvider(t *testing.T) {
	for _, e := range es {
		t.Run(e.Cloud.Name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					if e.Path == r.URL.Path {
						w.WriteHeader(200)
					} else {
						w.WriteHeader(404)
					}
				}))
			defer ts.Close()

			c := getCloudProvider(ts.URL)
			assert.Equal(t, e.Cloud.Name, c.Name, "Cloud name should match")
		})
	}
}
