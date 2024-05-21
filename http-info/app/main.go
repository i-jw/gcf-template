// Package helloworld provides a set of Cloud Functions samples.
package hello

import (
	"fmt"
	"io"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("entryPoint", entryPoint)
}

// HelloHTTP is an HTTP Cloud Function with a request parameter.
func entryPoint(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
	}
	resp, err := http.Get("http://ipinfo.io")
	if err != nil {
		fmt.Fprint(w, "outgoing ip request get err")
	}
	defer resp.Body.Close()
	body, b_err := io.ReadAll(resp.Body)
	if b_err == nil {
		fmt.Fprint(w, string(body))
		return
	}
}
